# Example FastAPI API for testing SNAG
# Installation: pip install fastapi uvicorn
# Run: uvicorn main:app --reload --port 8000

from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel
from typing import List, Optional

app = FastAPI(
    title="Example APIs for SNAG",
    description="A test API to demonstrate SNAG capabilities",
    version="1.0.0",
)


# Modelos
class User(BaseModel):
    id: Optional[int] = None
    name: str
    email: str
    age: Optional[int] = None


class Product(BaseModel):
    id: Optional[int] = None
    name: str
    price: float
    description: Optional[str] = None
    stock: int


class Order(BaseModel):
    id: Optional[int] = None
    user_id: int
    product_id: int
    quantity: int


class Category(BaseModel):
    id: Optional[int] = None
    name: str
    description: Optional[str] = None
    parent_id: Optional[int] = None


class Review(BaseModel):
    id: Optional[int] = None
    user_id: int
    product_id: int
    rating: int
    title: str
    body: Optional[str] = None


class Address(BaseModel):
    id: Optional[int] = None
    user_id: int
    street: str
    city: str
    state: str
    zip_code: str
    country: str
    is_default: Optional[bool] = False


class Coupon(BaseModel):
    id: Optional[int] = None
    code: str
    discount_percent: int
    max_uses: int
    active: Optional[bool] = True


class LoginRequest(BaseModel):
    email: str
    password: str


class Tag(BaseModel):
    id: Optional[int] = None
    name: str
    color: Optional[str] = "#ffffff"


# Base de datos simulada
users_db = [
    {"id": 1, "name": "Alice Johnson", "email": "alice@example.com", "age": 30},
    {"id": 2, "name": "Bob Smith", "email": "bob@example.com", "age": 25},
    {"id": 3, "name": "Carol White", "email": "carol@example.com", "age": 35},
]

products_db = [
    {"id": 1, "name": "Laptop", "price": 999.99, "description": "Powerful laptop", "stock": 10},
    {"id": 2, "name": "Mouse", "price": 29.99, "description": "Wireless mouse", "stock": 50},
    {"id": 3, "name": "Keyboard", "price": 49.99, "description": "RGB mechanical keyboard", "stock": 30},
    {"id": 4, "name": "Monitor", "price": 349.99, "description": "27-inch 4K monitor", "stock": 8},
]

orders_db = []

categories_db = [
    {"id": 1, "name": "Electronics", "description": "Electronic devices", "parent_id": None},
    {"id": 2, "name": "Peripherals", "description": "PC accessories", "parent_id": 1},
    {"id": 3, "name": "Audio", "description": "Audio equipment", "parent_id": 1},
]

reviews_db = [
    {"id": 1, "user_id": 1, "product_id": 1, "rating": 5, "title": "Excellent laptop", "body": "Very fast and durable."},
    {"id": 2, "user_id": 2, "product_id": 2, "rating": 4, "title": "Good mouse", "body": "Comfortable for long sessions."},
]

addresses_db = [
    {"id": 1, "user_id": 1, "street": "123 Main St", "city": "New York", "state": "NY", "zip_code": "10001", "country": "US", "is_default": True},
]

coupons_db = [
    {"id": 1, "code": "SNAG10", "discount_percent": 10, "max_uses": 100, "active": True},
    {"id": 2, "code": "SAVE20", "discount_percent": 20, "max_uses": 50, "active": True},
]

tags_db = [
    {"id": 1, "name": "new", "color": "#3fb950"},
    {"id": 2, "name": "sale", "color": "#f85149"},
    {"id": 3, "name": "featured", "color": "#bc8cff"},
]

# ================== USER ENDPOINTS ==================


@app.get("/users", tags=["Users"], summary="Get all users")
def get_users() -> List[User]:
    """Returns the full list of users"""
    return users_db


@app.get("/users/{user_id}", tags=["Users"], summary="Get user by ID")
def get_user(user_id: int) -> User:
    """Returns a specific user by their ID"""
    for user in users_db:
        if user["id"] == user_id:
            return user
    raise HTTPException(status_code=404, detail="User not found")


@app.post("/users", tags=["Users"], summary="Create new user")
def create_user(user: User) -> User:
    """Creates a new user in the system"""
    new_id = max([u["id"] for u in users_db]) + 1 if users_db else 1
    user_dict = user.dict()
    user_dict["id"] = new_id
    users_db.append(user_dict)
    return user_dict


@app.put("/users/{user_id}", tags=["Users"], summary="Update user")
def update_user(user_id: int, user: User) -> User:
    """Updates an existing user"""
    for i, u in enumerate(users_db):
        if u["id"] == user_id:
            user_dict = user.dict()
            user_dict["id"] = user_id
            users_db[i] = user_dict
            return user_dict
    raise HTTPException(status_code=404, detail="User not found")


@app.delete("/users/{user_id}", tags=["Users"], summary="Delete user")
def delete_user(user_id: int):
    """Deletes a user from the system"""
    for i, u in enumerate(users_db):
        if u["id"] == user_id:
            users_db.pop(i)
            return {"message": "User deleted successfully"}
    raise HTTPException(status_code=404, detail="User not found")


# ================== PRODUCT ENDPOINTS ==================


@app.get("/products", tags=["Products"], summary="Get all products")
def get_products() -> List[Product]:
    """Returns the full list of products"""
    return products_db


@app.get("/products/{product_id}", tags=["Products"], summary="Get product by ID")
def get_product(product_id: int) -> Product:
    """Returns a specific product by its ID"""
    for product in products_db:
        if product["id"] == product_id:
            return product
    raise HTTPException(status_code=404, detail="Product not found")


@app.post("/products", tags=["Products"], summary="Create new product")
def create_product(product: Product) -> Product:
    """Creates a new product in the inventory"""
    new_id = max([p["id"] for p in products_db]) + 1 if products_db else 1
    product_dict = product.dict()
    product_dict["id"] = new_id
    products_db.append(product_dict)
    return product_dict


@app.patch(
    "/products/{product_id}/stock", tags=["Products"], summary="Update stock"
)
def update_stock(product_id: int, quantity: int) -> Product:
    """Updates the stock quantity of a product"""
    for i, p in enumerate(products_db):
        if p["id"] == product_id:
            products_db[i]["stock"] = quantity
            return products_db[i]
    raise HTTPException(status_code=404, detail="Product not found")


# ================== ORDER ENDPOINTS ==================


@app.get("/orders", tags=["Orders"], summary="Get all orders")
def get_orders() -> List[Order]:
    """Returns the full list of orders"""
    return orders_db


@app.post("/orders", tags=["Orders"], summary="Create new order")
def create_order(order: Order) -> Order:
    """Creates a new purchase order"""
    # Verify the user exists
    user_exists = any(u["id"] == order.user_id for u in users_db)
    if not user_exists:
        raise HTTPException(status_code=404, detail="User not found")

    # Verify the product exists and has stock
    product = next((p for p in products_db if p["id"] == order.product_id), None)
    if not product:
        raise HTTPException(status_code=404, detail="Product not found")

    if product["stock"] < order.quantity:
        raise HTTPException(status_code=400, detail="Insufficient stock")

    # Create order
    new_id = max([o["id"] for o in orders_db]) + 1 if orders_db else 1
    order_dict = order.dict()
    order_dict["id"] = new_id
    orders_db.append(order_dict)

    # Update stock
    for p in products_db:
        if p["id"] == order.product_id:
            p["stock"] -= order.quantity
            break

    return order_dict


# ================== HEALTH ENDPOINT ==================


@app.get("/health", tags=["System"], summary="Health check")
def health_check():
    """Checks that the API is running correctly"""
    return {
        "status": "healthy",
        "version": "1.0.0",
        "users_count": len(users_db),
        "products_count": len(products_db),
        "orders_count": len(orders_db),
    }


# ================== CATEGORY ENDPOINTS ==================


@app.get("/categories", tags=["Categories"], summary="List categories")
def get_categories():
    """Returns all available categories"""
    return categories_db


@app.get("/categories/{category_id}", tags=["Categories"], summary="Get category by ID")
def get_category(category_id: int):
    """Returns a specific category"""
    for c in categories_db:
        if c["id"] == category_id:
            return c
    raise HTTPException(status_code=404, detail="Category not found")


@app.post("/categories", tags=["Categories"], summary="Create category")
def create_category(category: Category):
    """Creates a new category"""
    new_id = max([c["id"] for c in categories_db]) + 1 if categories_db else 1
    d = category.dict()
    d["id"] = new_id
    categories_db.append(d)
    return d


@app.put("/categories/{category_id}", tags=["Categories"], summary="Update category")
def update_category(category_id: int, category: Category):
    """Updates an existing category"""
    for i, c in enumerate(categories_db):
        if c["id"] == category_id:
            d = category.dict()
            d["id"] = category_id
            categories_db[i] = d
            return d
    raise HTTPException(status_code=404, detail="Category not found")


@app.delete("/categories/{category_id}", tags=["Categories"], summary="Delete category")
def delete_category(category_id: int):
    """Deletes a category"""
    for i, c in enumerate(categories_db):
        if c["id"] == category_id:
            categories_db.pop(i)
            return {"message": "Category deleted"}
    raise HTTPException(status_code=404, detail="Category not found")


# ================== REVIEW ENDPOINTS ==================


@app.get("/reviews", tags=["Reviews"], summary="List all reviews")
def get_reviews(product_id: Optional[int] = Query(None), user_id: Optional[int] = Query(None)):
    """Returns reviews with optional filters by product or user"""
    result = reviews_db
    if product_id is not None:
        result = [r for r in result if r["product_id"] == product_id]
    if user_id is not None:
        result = [r for r in result if r["user_id"] == user_id]
    return result


@app.get("/reviews/{review_id}", tags=["Reviews"], summary="Get review by ID")
def get_review(review_id: int):
    """Returns a specific review"""
    for r in reviews_db:
        if r["id"] == review_id:
            return r
    raise HTTPException(status_code=404, detail="Review not found")


@app.post("/reviews", tags=["Reviews"], summary="Create review")
def create_review(review: Review):
    """Creates a new review for a product"""
    if not any(u["id"] == review.user_id for u in users_db):
        raise HTTPException(status_code=404, detail="User not found")
    if not any(p["id"] == review.product_id for p in products_db):
        raise HTTPException(status_code=404, detail="Product not found")
    if not (1 <= review.rating <= 5):
        raise HTTPException(status_code=400, detail="Rating must be between 1 and 5")
    new_id = max([r["id"] for r in reviews_db]) + 1 if reviews_db else 1
    d = review.dict()
    d["id"] = new_id
    reviews_db.append(d)
    return d


@app.put("/reviews/{review_id}", tags=["Reviews"], summary="Update review")
def update_review(review_id: int, review: Review):
    """Updates an existing review"""
    for i, r in enumerate(reviews_db):
        if r["id"] == review_id:
            d = review.dict()
            d["id"] = review_id
            reviews_db[i] = d
            return d
    raise HTTPException(status_code=404, detail="Review not found")


@app.delete("/reviews/{review_id}", tags=["Reviews"], summary="Delete review")
def delete_review(review_id: int):
    """Deletes a review"""
    for i, r in enumerate(reviews_db):
        if r["id"] == review_id:
            reviews_db.pop(i)
            return {"message": "Review deleted"}
    raise HTTPException(status_code=404, detail="Review not found")


# ================== ADDRESS ENDPOINTS ==================


@app.get("/users/{user_id}/addresses", tags=["Addresses"], summary="Get addresses for a user")
def get_addresses(user_id: int):
    """Returns all addresses for a user"""
    if not any(u["id"] == user_id for u in users_db):
        raise HTTPException(status_code=404, detail="User not found")
    return [a for a in addresses_db if a["user_id"] == user_id]


@app.post("/users/{user_id}/addresses", tags=["Addresses"], summary="Add address")
def create_address(user_id: int, address: Address):
    """Adds a new address to a user"""
    if not any(u["id"] == user_id for u in users_db):
        raise HTTPException(status_code=404, detail="User not found")
    new_id = max([a["id"] for a in addresses_db]) + 1 if addresses_db else 1
    d = address.dict()
    d["id"] = new_id
    d["user_id"] = user_id
    addresses_db.append(d)
    return d


@app.put("/users/{user_id}/addresses/{address_id}", tags=["Addresses"], summary="Update address")
def update_address(user_id: int, address_id: int, address: Address):
    """Updates an existing address"""
    for i, a in enumerate(addresses_db):
        if a["id"] == address_id and a["user_id"] == user_id:
            d = address.dict()
            d["id"] = address_id
            d["user_id"] = user_id
            addresses_db[i] = d
            return d
    raise HTTPException(status_code=404, detail="Address not found")


@app.delete("/users/{user_id}/addresses/{address_id}", tags=["Addresses"], summary="Delete address")
def delete_address(user_id: int, address_id: int):
    """Deletes an address from a user"""
    for i, a in enumerate(addresses_db):
        if a["id"] == address_id and a["user_id"] == user_id:
            addresses_db.pop(i)
            return {"message": "Address deleted"}
    raise HTTPException(status_code=404, detail="Address not found")


# ================== COUPON ENDPOINTS ==================


@app.get("/coupons", tags=["Coupons"], summary="List coupons")
def get_coupons(active_only: bool = Query(False)):
    """Returns all coupons. Filter by active ones with ?active_only=true"""
    if active_only:
        return [c for c in coupons_db if c["active"]]
    return coupons_db


@app.get("/coupons/{code}", tags=["Coupons"], summary="Validate coupon by code")
def get_coupon(code: str):
    """Validates and returns a coupon by its code"""
    for c in coupons_db:
        if c["code"].upper() == code.upper():
            return c
    raise HTTPException(status_code=404, detail="Coupon not found")


@app.post("/coupons", tags=["Coupons"], summary="Create coupon")
def create_coupon(coupon: Coupon):
    """Creates a new discount coupon"""
    if any(c["code"].upper() == coupon.code.upper() for c in coupons_db):
        raise HTTPException(status_code=400, detail="A coupon with that code already exists")
    new_id = max([c["id"] for c in coupons_db]) + 1 if coupons_db else 1
    d = coupon.dict()
    d["id"] = new_id
    coupons_db.append(d)
    return d


@app.patch("/coupons/{coupon_id}/deactivate", tags=["Coupons"], summary="Deactivate coupon")
def deactivate_coupon(coupon_id: int):
    """Deactivates an existing coupon"""
    for c in coupons_db:
        if c["id"] == coupon_id:
            c["active"] = False
            return c
    raise HTTPException(status_code=404, detail="Coupon not found")


@app.delete("/coupons/{coupon_id}", tags=["Coupons"], summary="Delete coupon")
def delete_coupon(coupon_id: int):
    """Deletes a coupon"""
    for i, c in enumerate(coupons_db):
        if c["id"] == coupon_id:
            coupons_db.pop(i)
            return {"message": "Coupon deleted"}
    raise HTTPException(status_code=404, detail="Coupon not found")


# ================== TAG ENDPOINTS ==================


@app.get("/tags", tags=["Tags"], summary="List tags")
def get_tags():
    """Returns all available tags"""
    return tags_db


@app.post("/tags", tags=["Tags"], summary="Create tag")
def create_tag(tag: Tag):
    """Creates a new tag"""
    new_id = max([t["id"] for t in tags_db]) + 1 if tags_db else 1
    d = tag.dict()
    d["id"] = new_id
    tags_db.append(d)
    return d


@app.put("/tags/{tag_id}", tags=["Tags"], summary="Update tag")
def update_tag(tag_id: int, tag: Tag):
    """Updates an existing tag"""
    for i, t in enumerate(tags_db):
        if t["id"] == tag_id:
            d = tag.dict()
            d["id"] = tag_id
            tags_db[i] = d
            return d
    raise HTTPException(status_code=404, detail="Tag not found")


@app.delete("/tags/{tag_id}", tags=["Tags"], summary="Delete tag")
def delete_tag(tag_id: int):
    """Deletes a tag"""
    for i, t in enumerate(tags_db):
        if t["id"] == tag_id:
            tags_db.pop(i)
            return {"message": "Tag deleted"}
    raise HTTPException(status_code=404, detail="Tag not found")


# ================== AUTH ENDPOINTS ==================


@app.post("/auth/login", tags=["Auth"], summary="Login")
def login(credentials: LoginRequest):
    """Authenticates a user and returns a mock session token"""
    user = next((u for u in users_db if u["email"] == credentials.email), None)
    if not user:
        raise HTTPException(status_code=401, detail="Invalid credentials")
    return {
        "token": f"mock-token-{user['id']}-abc123",
        "user_id": user["id"],
        "name": user["name"],
        "email": user["email"],
    }


@app.post("/auth/logout", tags=["Auth"], summary="Logout")
def logout():
    """Invalidates the current session token"""
    return {"message": "Logged out successfully"}


@app.get("/auth/me", tags=["Auth"], summary="Authenticated user profile")
def me():
    """Returns the authenticated user profile (mock: always returns user 1)"""
    return users_db[0]


# ================== STATS ENDPOINTS (System) ==================


@app.get("/stats", tags=["System"], summary="General statistics")
def get_stats():
    """Returns global system statistics"""
    total_stock = sum(p["stock"] for p in products_db)
    avg_price = sum(p["price"] for p in products_db) / len(products_db) if products_db else 0
    return {
        "users": len(users_db),
        "products": len(products_db),
        "orders": len(orders_db),
        "categories": len(categories_db),
        "reviews": len(reviews_db),
        "coupons_active": len([c for c in coupons_db if c["active"]]),
        "total_stock": total_stock,
        "avg_product_price": round(avg_price, 2),
    }


@app.delete("/stats/reset", tags=["System"], summary="Reset test data")
def reset_data():
    """Clears orders and reviews to restart testing"""
    orders_db.clear()
    reviews_db.clear()
    return {"message": "Data reset", "orders": 0, "reviews": 0}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
