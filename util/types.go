package util

import "time"

type YAMLFile struct {
	Database Database `yaml:"database"`
	Secret   Secret   `yaml:"secret"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Secret struct {
	Hash  string `yaml:"hash"`
	Block string `yaml:"block"`
}

type Category struct {
	CategoryId    int
	CategoryName  string
	Subcategories []Subcategory
}

type Subcategory struct {
	SubCategoryId   int
	SubCategoryName string
	PermissionLevel int
}

type FullCategory struct {
	CategoryId    int
	CategoryName  string
	Subcategories []FullSubCategory
	Priority      int
}

type FullSubCategory struct {
	SubCategoryId   int
	SubCategoryName string
	PermissionLevel int
	Priority        int
}

type ThreadSummary struct {
	Threadid    int
	ThreadName  string
	Userid      int
	Username    string
	CreateDate  time.Time
	LastUpdated time.Time
}

type User struct {
	Userid          int
	Username        string
	CreationDate    time.Time
	PermissionLevel int
}

type ProfilePost struct {
	Postid       int
	Threadid     int
	ThreadName   string
	Userid       int
	Username     string
	PostText     string
	PostDate     time.Time
	PostEditDate *time.Time
}

type ThreadPost struct {
	PostId   int
	UserId   int
	Username string
	PostText string
	PostDate time.Time
	EditDate *time.Time
}

type Thread struct {
	Threadid   int
	ThreadName string
	UserId     int
	Username   string
	ThreadText string
	CreateDate time.Time
	EditDate   *time.Time
}
