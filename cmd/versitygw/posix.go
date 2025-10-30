// Copyright 2023 Versity Software
// This file is licensed under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/versity/versitygw/backend/meta"
	"github.com/versity/versitygw/backend/posix"
	"io/fs"
	"math"
	"path/filepath"
)

var (
	chownuid, chowngid bool
	bucketlinks        bool
	versioningDir      string
	dirPerms           uint
	sidecar            string
	nometa             bool
	forceNoTmpFile     bool
	baseUrl            string
	endpointPath       string
)

var skipdirs = cli.NewStringSlice("VisualDiscSpace", "IsoBuffer")

func posixCommand() *cli.Command {
	return &cli.Command{
		Name:  "posix",
		Usage: "posix filesystem storage backend",
		Description: `Any posix filesystem that supports extended attributes. The top level
directory for the gateway must be provided. All sub directories of the
top level directory are treated as buckets, and all files/directories
below the "bucket directory" are treated as the objects. The object
name is split on "/" separator to translate to posix storage.
For example:
top level: /mnt/fs/gwroot
bucket: mybucket
object: a/b/c/myobject
will be translated into the file /mnt/fs/gwroot/mybucket/a/b/c/myobject`,
		Action: runPosix,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "chuid",
				Usage:       "chown newly created files and directories to client account UID",
				EnvVars:     []string{"VGW_CHOWN_UID"},
				Destination: &chownuid,
			},
			&cli.BoolFlag{
				Name:        "chgid",
				Usage:       "chown newly created files and directories to client account GID",
				EnvVars:     []string{"VGW_CHOWN_GID"},
				Destination: &chowngid,
			},
			&cli.BoolFlag{
				Name:        "bucketlinks",
				Usage:       "allow symlinked directories at bucket level to be treated as buckets",
				EnvVars:     []string{"VGW_BUCKET_LINKS"},
				Destination: &bucketlinks,
			},
			&cli.StringFlag{
				Name:        "versioning-dir",
				Usage:       "the directory path to enable bucket versioning",
				EnvVars:     []string{"VGW_VERSIONING_DIR"},
				Destination: &versioningDir,
			},
			&cli.UintFlag{
				Name:        "dir-perms",
				Usage:       "default directory permissions for new directories",
				EnvVars:     []string{"VGW_DIR_PERMS"},
				Destination: &dirPerms,
				DefaultText: "0755",
				Value:       0755,
			},
			&cli.StringFlag{
				Name:        "sidecar",
				Usage:       "use provided sidecar directory to store metadata",
				EnvVars:     []string{"VGW_META_SIDECAR"},
				Destination: &sidecar,
			},
			&cli.StringFlag{
				Name:        "notify-base-url",
				Usage:       "Base URL of the notification service (e.g. http://192.168.78.213:8080)",
				EnvVars:     []string{"NOTIFY_BASE_URL"},
				Destination: &baseUrl,
			},
			&cli.StringFlag{
				Name:        "notify-endpoint-path",
				Usage:       "API endpoint path for creating notification tasks (e.g. /oss/rest/restServer/creatArcTask)",
				EnvVars:     []string{"NOTIFY_ENDPOINT_PATH"},
				Destination: &endpointPath,
			},
			&cli.StringSliceFlag{
				Name:    "skipdirs",
				Usage:   "Directories to skip (can specify multiple times)",
				EnvVars: []string{"SKIP_DIRS"},
				Value:   skipdirs,
			},
			&cli.BoolFlag{
				Name:        "nometa",
				Usage:       "disable metadata storage",
				EnvVars:     []string{"VGW_META_NONE"},
				Destination: &nometa,
			},
			&cli.BoolFlag{
				Name:        "disableotmp",
				Usage:       "disable O_TMPFILE support for new objects",
				EnvVars:     []string{"VGW_DISABLE_OTMP"},
				Destination: &forceNoTmpFile,
			},
		},
	}
}

func runPosix(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return fmt.Errorf("no directory provided for operation")
	}

	gwroot := (ctx.Args().Get(0))

	if dirPerms > math.MaxUint32 {
		return fmt.Errorf("invalid directory permissions: %d", dirPerms)
	}

	if nometa && sidecar != "" {
		return fmt.Errorf("cannot use both nometa and sidecar metadata")
	}

	var directorship []string
	for _, d := range skipdirs.Value() {
		directorship = append(directorship, d)
	}

	opts := posix.PosixOpts{
		ChownUID:       chownuid,
		ChownGID:       chowngid,
		BucketLinks:    bucketlinks,
		VersioningDir:  versioningDir,
		NewDirPerm:     fs.FileMode(dirPerms),
		ForceNoTmpFile: forceNoTmpFile,
		BaseUrl:        baseUrl,
		EndpointPath:   endpointPath,
		Skipdirs:       directorship,
	}

	var ms meta.MetadataStorer
	switch {
	case sidecar != "":
		sc, err := meta.NewSideCar(sidecar)
		if err != nil {
			return fmt.Errorf("failed to init sidecar metadata: %w", err)
		}
		ms = sc
		opts.SideCarDir = sidecar
	case nometa:
		ms = meta.NoMeta{}
	default:
		sm, err := meta.NewSqlMeta(filepath.Join(gwroot, "mate.sqlite"))
		if err != nil {
			return fmt.Errorf("failed to init sqlite metadata: %w", err)
		}
		ms = sm
	}

	be, err := posix.New(gwroot, ms, opts)
	if err != nil {
		return fmt.Errorf("failed to init posix backend: %w", err)
	}

	return runGateway(ctx.Context, be)
}
