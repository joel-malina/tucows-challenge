table "order_items" {
  schema = schema.public
  column "item_id" {
    null = false
    type = uuid
  }
  column "order_id" {
    null = false
    type = uuid
  }
  column "product_id" {
    null = false
    type = uuid
  }
  column "quantity" {
    null = false
    type = integer
  }
  column "price" {
    null = true
    type = numeric(10,2)
  }
  primary_key {
    columns = [column.item_id]
  }
  foreign_key "order_items_order_id_fkey" {
    columns     = [column.order_id]
    ref_columns = [table.orders.column.order_id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "order_items_product_id_fkey" {
    columns     = [column.product_id]
    ref_columns = [table.products.column.product_id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "orders" {
  schema = schema.public
  column "order_id" {
    null = false
    type = uuid
  }
  column "customer_id" {
    null = false
    type = uuid
  }
  column "order_date" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "status" {
    null = true
    type = character_varying(50)
  }
  column "total_price" {
    null = true
    type = numeric(10,2)
  }
  primary_key {
    columns = [column.order_id]
  }
}
table "products" {
  schema = schema.public
  column "product_id" {
    null = false
    type = uuid
  }
  column "name" {
    null = false
    type = character_varying(255)
  }
  column "description" {
    null = true
    type = text
  }
  column "price" {
    null = false
    type = numeric(10,2)
  }
  column "stock" {
    null = false
    type = integer
  }
  primary_key {
    columns = [column.product_id]
  }
}
schema "public" {
  comment = "standard public schema"
}
