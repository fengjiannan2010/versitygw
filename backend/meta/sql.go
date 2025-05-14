package meta

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

// SqlMeta is a metadata storer that uses SQLite + GORM to store metadata.
type SqlMeta struct {
	db *gorm.DB
}

// MetadataEntry represents a single metadata record.
type MetadataEntry struct {
	Bucket    string `gorm:"index:idx_bucket_object_attr,unique"`
	Object    string `gorm:"index:idx_bucket_object_attr,unique"`
	Attribute string `gorm:"index:idx_bucket_object_attr,unique"`
	Value     []byte
}

var (
	// ErrNoSuchAttributeKey is returned when the key does not exist.
	ErrNoSuchAttributeKey = errors.New("no such key")
)

// NewSqlMeta creates a new SqlMeta metadata storer using SQLite.
func NewSqlMeta(dbPath string) (SqlMeta, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return SqlMeta{}, fmt.Errorf("failed to open database: %v", err)
	}
	if err := db.AutoMigrate(&MetadataEntry{}); err != nil {
		return SqlMeta{}, fmt.Errorf("failed to migrate schema: %v", err)
	}
	return SqlMeta{db: db}, nil
}

// RetrieveAttribute gets a specific metadata value.
func (s SqlMeta) RetrieveAttribute(_ *os.File, bucket, object, attribute string) ([]byte, error) {
	var entry MetadataEntry
	res := s.db.Where("bucket = ? AND object = ? AND attribute = ?", bucket, object, attribute).First(&entry)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNoSuchAttributeKey
	}
	if res.Error != nil {
		return nil, fmt.Errorf("failed to retrieve attribute: %v", res.Error)
	}
	return entry.Value, nil
}

// StoreAttribute stores or updates a metadata value.
func (s SqlMeta) StoreAttribute(_ *os.File, bucket, object, attribute string, value []byte) error {
	entry := MetadataEntry{
		Bucket:    bucket,
		Object:    object,
		Attribute: attribute,
		Value:     value,
	}
	err := s.db.Where(MetadataEntry{Bucket: bucket, Object: object, Attribute: attribute}).Assign(MetadataEntry{
		Bucket:    bucket,
		Object:    object,
		Attribute: attribute,
		Value:     value,
	}).FirstOrCreate(&entry).Error
	return err
}

// DeleteAttribute removes a specific metadata attribute.
func (s SqlMeta) DeleteAttribute(bucket, object, attribute string) error {
	res := s.db.Where("bucket = ? AND object = ? AND attribute = ?", bucket, object, attribute).Delete(&MetadataEntry{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete attribute: %v", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNoSuchAttributeKey
	}
	return nil
}

// ListAttributes returns a list of attribute names for a bucket or object.
func (s SqlMeta) ListAttributes(bucket, object string) ([]string, error) {
	var entries []MetadataEntry
	res := s.db.Select("attribute").Where("bucket = ? AND object = ?", bucket, object).Find(&entries)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to list attributes: %v", res.Error)
	}
	var attrs []string
	for _, e := range entries {
		attrs = append(attrs, e.Attribute)
	}
	return attrs, nil
}

// DeleteAttributes removes all attributes for a bucket or object.
func (s SqlMeta) DeleteAttributes(bucket, object string) error {
	res := s.db.Where("bucket = ? AND object = ?", bucket, object).Delete(&MetadataEntry{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete attributes: %v", res.Error)
	}
	return nil
}
