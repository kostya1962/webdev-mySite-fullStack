package handlers

import (
	"myAPI/database"
	"myAPI/models"

	"github.com/gofiber/fiber/v2"
)

// BannerItem - структура, возвращаемая API
type BannerItem struct {
	ID        int            `json:"id"`
	ProductID int            `json:"product_id"`
	Image     string         `json:"image"`
	Position  int            `json:"position"`
	Product   *models.Product `json:"product,omitempty"`
}

func GetBanners(c *fiber.Ctx) error {
	// Return banner records with embedded product and category data
	query := `SELECT b.id, b.product_id, b.image, b.position,
		p.id, p.name, p.price, p.short_description, p.long_description,
		p.sku, p.discount, p.images, p.category_id, p.created_at, p.updated_at,
		c.id, c.name, c.alias
		FROM banners b
		JOIN products p ON b.product_id = p.id
		JOIN categories c ON p.category_id = c.id
		ORDER BY b.position ASC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch banners"})
	}
	defer rows.Close()

	var items []BannerItem
	for rows.Next() {
		var it BannerItem
		var prod models.Product
		var cat models.Category

		if err := rows.Scan(&it.ID, &it.ProductID, &it.Image, &it.Position,
			&prod.ID, &prod.Name, &prod.Price, &prod.ShortDescription, &prod.LongDescription,
			&prod.SKU, &prod.Discount, &prod.Images, &prod.CategoryID, &prod.CreatedAt, &prod.UpdatedAt,
			&cat.ID, &cat.Name, &cat.Alias);
			err != nil {
			continue
		}

		prod.Category = &cat
		it.Product = &prod
		items = append(items, it)
	}

	return c.JSON(fiber.Map{"banners": items})
}
