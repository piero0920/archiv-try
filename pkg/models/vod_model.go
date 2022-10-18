package models

import "time"

type Vod struct {
	UUID             string     `gorm:"primaryKey;uniqueIndex" json:"uuid" form:"uuid"`
	Title            string     `json:"title" form:"title"`
	Duration         int        `json:"duration" form:"duration"`
	Date             *time.Time `json:"date" form:"date" time_format:"2006-01-02T15:04:05.000Z"`
	Viewcount        int        `gorm:"default:0" json:"viewcount" form:"viewcount"`
	Filename         string     `json:"filename" form:"filename"`
	Resolution       string     `json:"resolution" form:"resolution"`
	Fps              float32    `json:"fps" form:"fps"`
	Size             int        `json:"size" form:"size"`
	Publish          bool       `json:"publish" form:"publish"`
	Clips            []Clip     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"clips" form:"clips"`
	Transcript       string     `gorm:"column:transcript;type:text" json:"transcript" form:"transcript"`
	TranscriptVector string     `gorm:"column:transcript_vector;->;type:tsvector generated always as (setweight(to_tsvector('german',transcript), 'A') || ' ' || setweight(to_tsvector('simple',transcript), 'B')) stored;default:(-);index:transcript_idx,type:gin" json:"-" form:"-"`
}
