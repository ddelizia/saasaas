package gormsaas

import (
	"context"
	"reflect"

	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/t"
	"gorm.io/gorm"
)

type Db struct {
	GormDb *gorm.DB
	db.Getter
	db.Creator
	db.Deleter
	db.Updater
	db.CursorLister
}

func (r *Db) Create(ctx context.Context, id t.String, data interface{}) error {
	reflect.ValueOf(data).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
	result := r.GormDb.WithContext(ctx).Create(data)
	return result.Error
}

func (r *Db) Get(ctx context.Context, id t.String, out interface{}) error {
	result := r.GormDb.WithContext(ctx).First(out, "id = ?", id)
	return result.Error
}

func (r *Db) Delete(ctx context.Context, id t.String, out interface{}) error {
	result := r.GormDb.WithContext(ctx).Delete(out, id)
	return result.Error
}

func (r *Db) Update(ctx context.Context, id t.String, data interface{}) error {
	getResult := t.GenerareEmptyInterface(data)
	err := r.Get(ctx, id, getResult)
	if err != nil {
		return err
	}

	result := r.GormDb.Model(t.GenerareEmptyInterface(data)).Where("id = ?", *id).Updates(data)
	return result.Error
}

func (r *Db) CursorList(ctx context.Context, startAtKey string, limit int64, out interface{}) (*db.CursorListResult, error) {
	return nil, nil
}
