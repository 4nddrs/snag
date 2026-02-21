# Ejemplo de API FastAPI para probar SNAG
# Instalación: pip install fastapi uvicorn
# Ejecutar: uvicorn example_api:app --reload

from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel
from typing import List, Optional

app = FastAPI(
    title="APIs de Ejemplo para SNAG",
    description="Una API de prueba para demostrar las capacidades de SNAG",
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
    {"id": 1, "name": "Juan Pérez", "email": "juan@example.com", "age": 30},
    {"id": 2, "name": "María García", "email": "maria@example.com", "age": 25},
    {"id": 3, "name": "Carlos López", "email": "carlos@example.com", "age": 35},
]

products_db = [
    {"id": 1, "name": "Laptop", "price": 999.99, "description": "Laptop potente", "stock": 10},
    {"id": 2, "name": "Mouse", "price": 29.99, "description": "Mouse inalámbrico", "stock": 50},
    {"id": 3, "name": "Teclado", "price": 49.99, "description": "Teclado mecánico RGB", "stock": 30},
    {"id": 4, "name": "Monitor", "price": 349.99, "description": "Monitor 4K 27 pulgadas", "stock": 8},
]

orders_db = []

categories_db = [
    {"id": 1, "name": "Electrónica", "description": "Dispositivos electrónicos", "parent_id": None},
    {"id": 2, "name": "Periféricos", "description": "Accesorios para PC", "parent_id": 1},
    {"id": 3, "name": "Audio", "description": "Equipos de audio", "parent_id": 1},
]

reviews_db = [
    {"id": 1, "user_id": 1, "product_id": 1, "rating": 5, "title": "Excelente laptop", "body": "Muy rápida y duradera."},
    {"id": 2, "user_id": 2, "product_id": 2, "rating": 4, "title": "Buen mouse", "body": "Cómodo para largas sesiones."},
]

addresses_db = [
    {"id": 1, "user_id": 1, "street": "Av. Reforma 123", "city": "CDMX", "state": "CDMX", "zip_code": "06600", "country": "MX", "is_default": True},
]

coupons_db = [
    {"id": 1, "code": "SNAG10", "discount_percent": 10, "max_uses": 100, "active": True},
    {"id": 2, "code": "SAVE20", "discount_percent": 20, "max_uses": 50, "active": True},
]

tags_db = [
    {"id": 1, "name": "nuevo", "color": "#3fb950"},
    {"id": 2, "name": "oferta", "color": "#f85149"},
    {"id": 3, "name": "destacado", "color": "#bc8cff"},
]

# ================== ENDPOINTS DE USUARIOS ==================


@app.get("/users", tags=["Users"], summary="Obtener todos los usuarios")
def get_users() -> List[User]:
    """Devuelve la lista completa de usuarios"""
    return users_db


@app.get("/users/{user_id}", tags=["Users"], summary="Obtener usuario por ID")
def get_user(user_id: int) -> User:
    """Devuelve un usuario específico por su ID"""
    for user in users_db:
        if user["id"] == user_id:
            return user
    raise HTTPException(status_code=404, detail="Usuario no encontrado")


@app.post("/users", tags=["Users"], summary="Crear nuevo usuario")
def create_user(user: User) -> User:
    """Crea un nuevo usuario en el sistema"""
    new_id = max([u["id"] for u in users_db]) + 1 if users_db else 1
    user_dict = user.dict()
    user_dict["id"] = new_id
    users_db.append(user_dict)
    return user_dict


@app.put("/users/{user_id}", tags=["Users"], summary="Actualizar usuario")
def update_user(user_id: int, user: User) -> User:
    """Actualiza un usuario existente"""
    for i, u in enumerate(users_db):
        if u["id"] == user_id:
            user_dict = user.dict()
            user_dict["id"] = user_id
            users_db[i] = user_dict
            return user_dict
    raise HTTPException(status_code=404, detail="Usuario no encontrado")


@app.delete("/users/{user_id}", tags=["Users"], summary="Eliminar usuario")
def delete_user(user_id: int):
    """Elimina un usuario del sistema"""
    for i, u in enumerate(users_db):
        if u["id"] == user_id:
            users_db.pop(i)
            return {"message": "Usuario eliminado exitosamente"}
    raise HTTPException(status_code=404, detail="Usuario no encontrado")


# ================== ENDPOINTS DE PRODUCTOS ==================


@app.get("/products", tags=["Products"], summary="Obtener todos los productos")
def get_products() -> List[Product]:
    """Devuelve la lista completa de productos"""
    return products_db


@app.get("/products/{product_id}", tags=["Products"], summary="Obtener producto por ID")
def get_product(product_id: int) -> Product:
    """Devuelve un producto específico por su ID"""
    for product in products_db:
        if product["id"] == product_id:
            return product
    raise HTTPException(status_code=404, detail="Producto no encontrado")


@app.post("/products", tags=["Products"], summary="Crear nuevo producto")
def create_product(product: Product) -> Product:
    """Crea un nuevo producto en el inventario"""
    new_id = max([p["id"] for p in products_db]) + 1 if products_db else 1
    product_dict = product.dict()
    product_dict["id"] = new_id
    products_db.append(product_dict)
    return product_dict


@app.patch(
    "/products/{product_id}/stock", tags=["Products"], summary="Actualizar stock"
)
def update_stock(product_id: int, quantity: int) -> Product:
    """Actualiza el stock de un producto"""
    for i, p in enumerate(products_db):
        if p["id"] == product_id:
            products_db[i]["stock"] = quantity
            return products_db[i]
    raise HTTPException(status_code=404, detail="Producto no encontrado")


# ================== ENDPOINTS DE ÓRDENES ==================


@app.get("/orders", tags=["Orders"], summary="Obtener todas las órdenes")
def get_orders() -> List[Order]:
    """Devuelve la lista completa de órdenes"""
    return orders_db


@app.post("/orders", tags=["Orders"], summary="Crear nueva orden")
def create_order(order: Order) -> Order:
    """Crea una nueva orden de compra"""
    # Verificar que el usuario existe
    user_exists = any(u["id"] == order.user_id for u in users_db)
    if not user_exists:
        raise HTTPException(status_code=404, detail="Usuario no encontrado")

    # Verificar que el producto existe y hay stock
    product = next((p for p in products_db if p["id"] == order.product_id), None)
    if not product:
        raise HTTPException(status_code=404, detail="Producto no encontrado")

    if product["stock"] < order.quantity:
        raise HTTPException(status_code=400, detail="Stock insuficiente")

    # Crear orden
    new_id = max([o["id"] for o in orders_db]) + 1 if orders_db else 1
    order_dict = order.dict()
    order_dict["id"] = new_id
    orders_db.append(order_dict)

    # Actualizar stock
    for p in products_db:
        if p["id"] == order.product_id:
            p["stock"] -= order.quantity
            break

    return order_dict


# ================== ENDPOINT DE SALUD ==================


@app.get("/health", tags=["System"], summary="Health check")
def health_check():
    """Verifica que la API está funcionando correctamente"""
    return {
        "status": "healthy",
        "version": "1.0.0",
        "users_count": len(users_db),
        "products_count": len(products_db),
        "orders_count": len(orders_db),
    }


# ================== ENDPOINTS DE CATEGORÍAS ==================


@app.get("/categories", tags=["Categories"], summary="Listar categorías")
def get_categories():
    """Devuelve todas las categorías disponibles"""
    return categories_db


@app.get("/categories/{category_id}", tags=["Categories"], summary="Obtener categoría por ID")
def get_category(category_id: int):
    """Devuelve una categoría específica"""
    for c in categories_db:
        if c["id"] == category_id:
            return c
    raise HTTPException(status_code=404, detail="Categoría no encontrada")


@app.post("/categories", tags=["Categories"], summary="Crear categoría")
def create_category(category: Category):
    """Crea una nueva categoría"""
    new_id = max([c["id"] for c in categories_db]) + 1 if categories_db else 1
    d = category.dict()
    d["id"] = new_id
    categories_db.append(d)
    return d


@app.put("/categories/{category_id}", tags=["Categories"], summary="Actualizar categoría")
def update_category(category_id: int, category: Category):
    """Actualiza una categoría existente"""
    for i, c in enumerate(categories_db):
        if c["id"] == category_id:
            d = category.dict()
            d["id"] = category_id
            categories_db[i] = d
            return d
    raise HTTPException(status_code=404, detail="Categoría no encontrada")


@app.delete("/categories/{category_id}", tags=["Categories"], summary="Eliminar categoría")
def delete_category(category_id: int):
    """Elimina una categoría"""
    for i, c in enumerate(categories_db):
        if c["id"] == category_id:
            categories_db.pop(i)
            return {"message": "Categoría eliminada"}
    raise HTTPException(status_code=404, detail="Categoría no encontrada")


# ================== ENDPOINTS DE REVIEWS ==================


@app.get("/reviews", tags=["Reviews"], summary="Listar todas las reviews")
def get_reviews(product_id: Optional[int] = Query(None), user_id: Optional[int] = Query(None)):
    """Devuelve reviews, con filtros opcionales por producto o usuario"""
    result = reviews_db
    if product_id is not None:
        result = [r for r in result if r["product_id"] == product_id]
    if user_id is not None:
        result = [r for r in result if r["user_id"] == user_id]
    return result


@app.get("/reviews/{review_id}", tags=["Reviews"], summary="Obtener review por ID")
def get_review(review_id: int):
    """Devuelve una review específica"""
    for r in reviews_db:
        if r["id"] == review_id:
            return r
    raise HTTPException(status_code=404, detail="Review no encontrada")


@app.post("/reviews", tags=["Reviews"], summary="Crear review")
def create_review(review: Review):
    """Crea una nueva review para un producto"""
    if not any(u["id"] == review.user_id for u in users_db):
        raise HTTPException(status_code=404, detail="Usuario no encontrado")
    if not any(p["id"] == review.product_id for p in products_db):
        raise HTTPException(status_code=404, detail="Producto no encontrado")
    if not (1 <= review.rating <= 5):
        raise HTTPException(status_code=400, detail="Rating debe estar entre 1 y 5")
    new_id = max([r["id"] for r in reviews_db]) + 1 if reviews_db else 1
    d = review.dict()
    d["id"] = new_id
    reviews_db.append(d)
    return d


@app.put("/reviews/{review_id}", tags=["Reviews"], summary="Editar review")
def update_review(review_id: int, review: Review):
    """Actualiza una review existente"""
    for i, r in enumerate(reviews_db):
        if r["id"] == review_id:
            d = review.dict()
            d["id"] = review_id
            reviews_db[i] = d
            return d
    raise HTTPException(status_code=404, detail="Review no encontrada")


@app.delete("/reviews/{review_id}", tags=["Reviews"], summary="Eliminar review")
def delete_review(review_id: int):
    """Elimina una review"""
    for i, r in enumerate(reviews_db):
        if r["id"] == review_id:
            reviews_db.pop(i)
            return {"message": "Review eliminada"}
    raise HTTPException(status_code=404, detail="Review no encontrada")


# ================== ENDPOINTS DE DIRECCIONES ==================


@app.get("/users/{user_id}/addresses", tags=["Addresses"], summary="Direcciones de un usuario")
def get_addresses(user_id: int):
    """Devuelve todas las direcciones de un usuario"""
    if not any(u["id"] == user_id for u in users_db):
        raise HTTPException(status_code=404, detail="Usuario no encontrado")
    return [a for a in addresses_db if a["user_id"] == user_id]


@app.post("/users/{user_id}/addresses", tags=["Addresses"], summary="Agregar dirección")
def create_address(user_id: int, address: Address):
    """Agrega una nueva dirección a un usuario"""
    if not any(u["id"] == user_id for u in users_db):
        raise HTTPException(status_code=404, detail="Usuario no encontrado")
    new_id = max([a["id"] for a in addresses_db]) + 1 if addresses_db else 1
    d = address.dict()
    d["id"] = new_id
    d["user_id"] = user_id
    addresses_db.append(d)
    return d


@app.put("/users/{user_id}/addresses/{address_id}", tags=["Addresses"], summary="Actualizar dirección")
def update_address(user_id: int, address_id: int, address: Address):
    """Actualiza una dirección existente"""
    for i, a in enumerate(addresses_db):
        if a["id"] == address_id and a["user_id"] == user_id:
            d = address.dict()
            d["id"] = address_id
            d["user_id"] = user_id
            addresses_db[i] = d
            return d
    raise HTTPException(status_code=404, detail="Dirección no encontrada")


@app.delete("/users/{user_id}/addresses/{address_id}", tags=["Addresses"], summary="Eliminar dirección")
def delete_address(user_id: int, address_id: int):
    """Elimina una dirección de un usuario"""
    for i, a in enumerate(addresses_db):
        if a["id"] == address_id and a["user_id"] == user_id:
            addresses_db.pop(i)
            return {"message": "Dirección eliminada"}
    raise HTTPException(status_code=404, detail="Dirección no encontrada")


# ================== ENDPOINTS DE CUPONES ==================


@app.get("/coupons", tags=["Coupons"], summary="Listar cupones")
def get_coupons(active_only: bool = Query(False)):
    """Devuelve todos los cupones. Filtra por activos con ?active_only=true"""
    if active_only:
        return [c for c in coupons_db if c["active"]]
    return coupons_db


@app.get("/coupons/{code}", tags=["Coupons"], summary="Validar cupón por código")
def get_coupon(code: str):
    """Valida y devuelve un cupón por su código"""
    for c in coupons_db:
        if c["code"].upper() == code.upper():
            return c
    raise HTTPException(status_code=404, detail="Cupón no encontrado")


@app.post("/coupons", tags=["Coupons"], summary="Crear cupón")
def create_coupon(coupon: Coupon):
    """Crea un nuevo cupón de descuento"""
    if any(c["code"].upper() == coupon.code.upper() for c in coupons_db):
        raise HTTPException(status_code=400, detail="Ya existe un cupón con ese código")
    new_id = max([c["id"] for c in coupons_db]) + 1 if coupons_db else 1
    d = coupon.dict()
    d["id"] = new_id
    coupons_db.append(d)
    return d


@app.patch("/coupons/{coupon_id}/deactivate", tags=["Coupons"], summary="Desactivar cupón")
def deactivate_coupon(coupon_id: int):
    """Desactiva un cupón existente"""
    for c in coupons_db:
        if c["id"] == coupon_id:
            c["active"] = False
            return c
    raise HTTPException(status_code=404, detail="Cupón no encontrado")


@app.delete("/coupons/{coupon_id}", tags=["Coupons"], summary="Eliminar cupón")
def delete_coupon(coupon_id: int):
    """Elimina un cupón"""
    for i, c in enumerate(coupons_db):
        if c["id"] == coupon_id:
            coupons_db.pop(i)
            return {"message": "Cupón eliminado"}
    raise HTTPException(status_code=404, detail="Cupón no encontrado")


# ================== ENDPOINTS DE TAGS ==================


@app.get("/tags", tags=["Tags"], summary="Listar tags")
def get_tags():
    """Devuelve todos los tags disponibles"""
    return tags_db


@app.post("/tags", tags=["Tags"], summary="Crear tag")
def create_tag(tag: Tag):
    """Crea un nuevo tag"""
    new_id = max([t["id"] for t in tags_db]) + 1 if tags_db else 1
    d = tag.dict()
    d["id"] = new_id
    tags_db.append(d)
    return d


@app.put("/tags/{tag_id}", tags=["Tags"], summary="Actualizar tag")
def update_tag(tag_id: int, tag: Tag):
    """Actualiza un tag existente"""
    for i, t in enumerate(tags_db):
        if t["id"] == tag_id:
            d = tag.dict()
            d["id"] = tag_id
            tags_db[i] = d
            return d
    raise HTTPException(status_code=404, detail="Tag no encontrado")


@app.delete("/tags/{tag_id}", tags=["Tags"], summary="Eliminar tag")
def delete_tag(tag_id: int):
    """Elimina un tag"""
    for i, t in enumerate(tags_db):
        if t["id"] == tag_id:
            tags_db.pop(i)
            return {"message": "Tag eliminado"}
    raise HTTPException(status_code=404, detail="Tag no encontrado")


# ================== ENDPOINTS DE AUTH ==================


@app.post("/auth/login", tags=["Auth"], summary="Iniciar sesión")
def login(credentials: LoginRequest):
    """Autentica un usuario y devuelve un token de sesión simulado"""
    user = next((u for u in users_db if u["email"] == credentials.email), None)
    if not user:
        raise HTTPException(status_code=401, detail="Credenciales inválidas")
    return {
        "token": f"mock-token-{user['id']}-abc123",
        "user_id": user["id"],
        "name": user["name"],
        "email": user["email"],
    }


@app.post("/auth/logout", tags=["Auth"], summary="Cerrar sesión")
def logout():
    """Invalida el token de sesión actual"""
    return {"message": "Sesión cerrada exitosamente"}


@app.get("/auth/me", tags=["Auth"], summary="Perfil del usuario autenticado")
def me():
    """Devuelve el perfil del usuario autenticado (simulado: siempre user 1)"""
    return users_db[0]


# ================== ENDPOINTS DE STATS (System) ==================


@app.get("/stats", tags=["System"], summary="Estadísticas generales")
def get_stats():
    """Devuelve estadísticas globales del sistema"""
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


@app.delete("/stats/reset", tags=["System"], summary="Resetear datos de prueba")
def reset_data():
    """Limpia órdenes y reviews para reiniciar las pruebas"""
    orders_db.clear()
    reviews_db.clear()
    return {"message": "Datos reseteados", "orders": 0, "reviews": 0}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
