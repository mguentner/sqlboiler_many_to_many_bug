package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	_ "github.com/lib/pq"

	"github.com/mguentner/sqlboiler_many_to_many_bug/models"
)

func TestNestedEagerLoadingManyToMany(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://postgres:secret@127.0.0.1:5432/sqlboiler_bug?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	db.ExecContext(ctx, "DELETE FROM list_items")
	db.ExecContext(ctx, "DELETE FROM item_tags")
	db.ExecContext(ctx, "DELETE FROM lists")
	db.ExecContext(ctx, "DELETE FROM items")
	db.ExecContext(ctx, "DELETE FROM tags")

	// Create a tag
	tag := &models.Tag{ID: "tag1", Name: "Important"}
	err = tag.Insert(ctx, db, boil.Infer())
	if err != nil {
		t.Fatal(err)
	}

	// Create an item
	item := &models.Item{ID: "item1", Name: "Test Item"}
	err = item.Insert(ctx, db, boil.Infer())
	if err != nil {
		t.Fatal(err)
	}

	// Add tag to item
	_, err = db.ExecContext(ctx, "INSERT INTO item_tags (item_id, tag_id) VALUES ($1, $2)", "item1", "tag1")
	if err != nil {
		t.Fatal(err)
	}

	// Let's add two lists
	list1 := &models.List{ID: "list1", Name: "List One"}
	err = list1.Insert(ctx, db, boil.Infer())
	if err != nil {
		t.Fatal(err)
	}

	list2 := &models.List{ID: "list2", Name: "List Two"}
	err = list2.Insert(ctx, db, boil.Infer())
	if err != nil {
		t.Fatal(err)
	}

	// The item is part of both lists:
	//      tag1
	//       |
	//      item1
	//     /     \
	//    list1 list2
	_, err = db.ExecContext(ctx, "INSERT INTO list_items (list_id, item_id) VALUES ($1, $2)", "list1", "item1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ExecContext(ctx, "INSERT INTO list_items (list_id, item_id) VALUES ($1, $2)", "list2", "item1")
	if err != nil {
		t.Fatal(err)
	}

	// Query all lists with nested eager loading:
	// Lists -> Items -> Tags (via many-to-many)
	lists, err := models.Lists(
		qm.Load(models.ListRels.Items),
		qm.Load(qm.Rels(models.ListRels.Items, models.ItemRels.Tags)),
	).All(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	if len(lists) != 2 {
		t.Fatalf("expected 2 lists, got %d", len(lists))
	}

	// BUG: One list's item will have Tags populated, the other will have empty slice
	for _, list := range lists {
		if len(list.R.Items) != 1 {
			t.Fatalf("list %s: expected 1 item, got %d", list.ID, len(list.R.Items))
		}

		item := list.R.Items[0]

		// THIS ASSERTION WILL FAIL FOR ONE OF THE LISTS
		if len(item.R.Tags) != 1 {
			t.Errorf("BUG: list %s, item %s: expected 1 Tag, got %d (nested relation not populated on duplicate instance)",
				list.ID, item.ID, len(item.R.Tags))
		}
	}
}
