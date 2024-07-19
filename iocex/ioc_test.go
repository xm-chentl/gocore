package iocex

import (
	"testing"
)

type IDb interface {
	Db() string
}

type dbInst struct{}

func (dbInst) Db() string {
	return "db"
}

type mongoDbInst struct{}

func (m mongoDbInst) Db() string {
	return "mongo"
}

type nestedInst struct {
	*APIInst
}

type APIInst struct {
	Db IDb `inject:""`
}

type apiInst struct {
	Db   IDb   `inject:""`
	Db2  IDb   `inject:"mongo"`
	Tool *Tool `inject:""`
}

type Tool struct {
	Name string
}

func Test_Get(t *testing.T) {
	key := getType(new(IDb))
	container.Store(
		key,
		&dbInst{},
	)
	v := Get(new(IDb))
	if v == nil {
		t.Fatal("err")
	}

	inst, ok := v.(IDb)
	if !ok {
		t.Fatal("err")
	}
	if inst.Db() != "db" {
		t.Fatal("err")
	}

	container.Delete(key)
}

func Test_GetKey(t *testing.T) {
	key := getType(new(IDb))
	container.Store(
		key,
		map[string]interface{}{
			"mongo": &mongoDbInst{},
		},
	)
	v := getTag(new(IDb), "mongo")
	if v == nil {
		t.Fatal("err")
	}

	inst, ok := v.(IDb)
	if !ok {
		t.Fatal("err")
	}
	if inst.Db() != "mongo" {
		t.Fatal("err")
	}

	container.Delete(key)
}

func Test_Has(t *testing.T) {
	key := getType(new(IDb))
	container.Store(
		key,
		&dbInst{},
	)
	if ok := Has(new(IDb)); !ok {
		t.Fatal("err")
	}
	container.Delete(key)
}

func Test_Injectss(t *testing.T) {
	key := getType(new(IDb))
	container.Store(
		key,
		&dbInst{},
	)
	key1 := getType(new(Tool))
	container.Store(
		key1,
		&Tool{
			Name: "test_001",
		},
	)
	inst := &apiInst{}
	if err := Inject(inst); err != nil {
		t.Fatal("err")
	}
	if inst.Db.Db() != "db" {
		t.Fatal("err")
	}
	if inst.Tool == nil {
		t.Fatal("err")
	}
	container.Delete(key)
}

func Test_Inject(t *testing.T) {

	t.Run(`Tag:inject:""`, func(test *testing.T) {
		key := getType(new(IDb))
		container.Store(
			key,
			&dbInst{},
		)
		key1 := getType(new(Tool))
		container.Store(
			key1,
			&Tool{
				Name: "test_001",
			},
		)
		inst := &apiInst{}
		if err := Inject(inst); err != nil {
			test.Fatal("err")
		}
		if inst.Db.Db() != "db" {
			test.Fatal("err")
		}
		if inst.Tool == nil {
			test.Fatal("err")
		}
		container.Delete(key)
	})

	t.Run(`Tag:inject:"mongo"`, func(test *testing.T) {
		key := getType(new(IDb))
		container.Store(
			key,
			map[string]interface{}{
				TagDefault: &dbInst{},
				"mongo":    &mongoDbInst{},
			},
		)
		inst := &apiInst{}
		if err := Inject(inst); err != nil {
			test.Fatal("err")
		}
		if inst.Db == nil {
			test.Fatal("err")
		}
		if inst.Db2.Db() != "mongo" {
			test.Fatal("err")
		}
		container.Delete(key)
	})

	t.Run(`nested inject`, func(test *testing.T) {
		key := getType(new(IDb))
		container.Store(
			key,
			&dbInst{},
		)

		inst := &nestedInst{}
		if err := Inject(inst); err != nil {
			test.Fatal("err")
		}
		if inst.Db == nil {
			test.Fatal("err")
		}
		if inst.Db.Db() != "db" {
			test.Fatal("err")
		}
		container.Delete(key)
	})
}

func Test_Set(t *testing.T) {
	key := new(IDb)
	Set(key, &dbInst{})
	v, ok := container.Load(getType(key))
	if !ok {
		t.Fatal("err")
	}

	db := v.(IDb)
	if db.Db() != "db" {
		t.Fatal("err")
	}
	container.Delete(key)
}

func Test_SetMap(t *testing.T) {
	key := new(IDb)
	SetMap(key, map[string]interface{}{
		"mongo": &mongoDbInst{},
	})

	v, ok := container.Load(getType(key))
	if !ok {
		t.Fatal("err")
	}

	vMap, ok := v.(map[string]interface{})
	if !ok {
		t.Fatal("err")
	}

	instV, ok := vMap["mongo"]
	if !ok {
		t.Fatal("err")
	}

	inst, ok := instV.(IDb)
	if !ok || inst.Db() != "mongo" {
		t.Fatal("err")
	}
	container.Delete(key)
}
