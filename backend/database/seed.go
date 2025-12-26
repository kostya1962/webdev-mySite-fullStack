package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	// Добавляем категории
	categories := []struct {
		name  string
		alias string
	}{
		{"Серьги", "earrings"},
		{"Кольца", "rings"},
		{"Ожерелья", "necklaces"},
		{"Заколки", "bracelets"},
	}

	for _, cat := range categories {
		_, err := DB.Exec("INSERT OR IGNORE INTO categories (name, alias) VALUES (?, ?)", cat.name, cat.alias)
		if err != nil {
			log.Printf("Failed to insert category %s: %v", cat.name, err)
		}
	}

	// Добавляем товары
	products := []struct {
		name             string
		price            float64
		shortDescription string
		longDescription  string
		sku              string
		discount         int
		images           []string
		categoryID       int
	}{
		{
			name:             "Lira Earrings",
			price:            1540.00,
			shortDescription: "Элегантные золотистые серьги-кольца",
			longDescription:  "Отлично подойдут к любому гардеробу. Чистое золото высокой пробы, которое не оставит вас равнодушными к качеству изделия.",
			sku:              "12",
			discount:         0,
			images:           []string{"/images/jewelry/lira1.jpg", "/images/jewelry/lira2.jpg", "/images/jewelry/lira3.jpg", "/images/jewelry/lira4.jpg"},
			categoryID:       1,
		},
		{
			name:             "Stella Diamond Ring",
			price:            92400.00,
			shortDescription: "Обручальное кольцо с бриллиантом",
			longDescription:  "Изысканное кольцо из белого золота 585 пробы с центральным бриллиантом весом 0.5 карата. Идеально для предложения руки и сердца.",
			sku:              "RING-STELLA-001",
			discount:         10,
			images:           []string{"/images/jewelry/stella1.jpg", "/images/jewelry/stella2.jpg"},
			categoryID:       2,
		},
		{
			name:             "Moonlight Necklace",
			price:            34650.00,
			shortDescription: "Изящное ожерелье с лунным камнем",
			longDescription:  "Утонченное ожерелье из серебра 925 пробы с натуральным лунным камнем. Подчеркнет вашу женственность и загадочность.",
			sku:              "NECK-MOON-001",
			discount:         15,
			images:           []string{"/images/jewelry/moonlight1.jpg", "/images/jewelry/moonlight2.jpg", "/images/jewelry/moonlight3.jpg"},
			categoryID:       3,
		},
		{
			name:             "Rose Gold Bracelet",
			price:            52360.00,
			shortDescription: "Браслет из розового золота",
			longDescription:  "Элегантный браслет из розового золота 750 пробы с тонким плетением. Идеальное дополнение к вечернему образу.",
			sku:              "BRACE-ROSE-001",
			discount:         0,
			images:           []string{"/images/jewelry/rosegold1.jpg"},
			categoryID:       4,
		},
		{
			name:             "Crystal Drop Earrings",
			price:            6545.00,
			shortDescription: "Серьги-капли с кристаллами Swarovski",
			longDescription:  "Потрясающие серьги с кристаллами Swarovski в оправе из родированного серебра. Добавят блеска любому вечернему наряду.",
			sku:              "EARR-CRYSTAL-001",
			discount:         20,
			images:           []string{"/images/jewelry/crystal1.jpg", "/images/jewelry/crystal2.jpg"},
			categoryID:       1,
		},
		{
			name:             "Vintage Pearl Necklace",
			price:            24640.00,
			shortDescription: "Винтажное жемчужное ожерелье",
			longDescription:  "Классическое ожерелье из натурального пресноводного жемчуга с позолоченной застежкой. Воплощение элегантности и изысканности.",
			sku:              "NECK-PEARL-001",
			discount:         0,
			images:           []string{"/images/jewelry/pearl1.jpg", "/images/jewelry/pearl2.jpg", "/images/jewelry/pearl3.jpg"},
			categoryID:       3,
		},
		{
			name:             "Infinity Ring Set",
			price:            11550.00,
			shortDescription: "Набор колец 'Бесконечность'",
			longDescription:  "Стильный набор из трех колец разного размера с символом бесконечности. Выполнен из серебра 925 пробы с родиевым покрытием.",
			sku:              "RING-INF-SET",
			discount:         25,
			images:           []string{"/images/jewelry/infinity1.jpg", "/images/jewelry/infinity2.jpg"},
			categoryID:       2,
		},
		{
			name:             "Charm Bracelet",
			price:            7315.00,
			shortDescription: "Браслет с подвесками-шармами",
			longDescription:  "Игривый браслет из серебра с коллекцией миниатюрных подвесок: сердечко, звездочка, луна и солнце. Отличный подарок для молодых девушек.",
			sku:              "BRACE-CHARM-001",
			discount:         0,
			images:           []string{"/images/jewelry/charm1.jpg", "/images/jewelry/charm2.jpg", "/images/jewelry/charm3.jpg"},
			categoryID:       4,
		},
	}

	for _, prod := range products {
		imagesJSON, _ := json.Marshal(prod.images)
		_, err := DB.Exec(`
			INSERT OR IGNORE INTO products 
			(name, price, short_description, long_description, sku, discount, images, category_id) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			prod.name, prod.price, prod.shortDescription, prod.longDescription,
			prod.sku, prod.discount, string(imagesJSON), prod.categoryID)
		if err != nil {
			log.Printf("Failed to insert product %s: %v", prod.name, err)
		}
	}

	// Добавляем отзывы
	reviews := []struct {
		productID int
		name      string
		text      string
		rating    int
	}{
		{1, "Василий", "В целом отлично, так как стоят не дорого и не пришлось жене дарить новый телефон, который она просила.", 3},
		{1, "Николай", "Я бы поставил и больше, но при доставке данные серёжки помялись и стали треугольными, а не круглыми. В остальном сервис отличный и я буду продолжать покупать в этом магазине.", 3},
		{2, "Светлана", "Превосходное кольцо! Муж сделал предложение именно с ним. Качество бриллианта на высоте!", 5},
		{2, "Андрей", "Дорого, но того стоит. Невеста в восторге!", 4},
		{3, "Елена", "Очень красивое ожерелье, лунный камень переливается на свету", 5},
		{4, "Мария", "Изысканный браслет, отлично подходит к моему стилю", 5},
		{5, "Анна", "Серьги просто великолепны! Кристаллы действительно сверкают", 5},
		{5, "Дарья", "Качественные серьги, но упаковка могла бы быть лучше", 4},
		{6, "Ольга", "Классическое жемчужное ожерелье, как у бабушки. Очень довольна!", 5},
		{7, "Юлия", "Набор колец супер! Ношу все сразу или по отдельности", 4},
		{8, "Екатерина", "Милый браслет, дочка в восторге от подвесок", 5},
	}

	for _, rev := range reviews {
		_, err := DB.Exec(`
			INSERT OR IGNORE INTO reviews (product_id, name, text, rating) 
			VALUES (?, ?, ?, ?)`,
			rev.productID, rev.name, rev.text, rev.rating)
		if err != nil {
			log.Printf("Failed to insert review: %v", err)
		}
	}

	// Добавляем записи для баннера (связываем с существующими товарами)
	banners := []struct{
		productID int
		image string
		position int
	}{
		{1, "/images/banner/bun1.png", 1},
		{3, "/images/banner/bun2.png", 2},
		{4, "/images/banner/bun3.png", 3},
	}

	for _, b := range banners {
		_, err := DB.Exec(`INSERT OR IGNORE INTO banners (product_id, image, position) VALUES (?, ?, ?)`, b.productID, b.image, b.position)
		if err != nil {
			log.Printf("Failed to insert banner item for product %d: %v", b.productID, err)
		}
	}

	// Добавляем новости (по умолчанию две записи)
	news := []struct{
		title string
		description string
		image string
	}{
		{
			title: "Торжественное открытие нашего ювелирного магазина",
			description: "Мы рады сообщить об открытии нового ювелирного магазина! В нашем ассортименте — изысканные украшения из золота и серебра, коллекции с драгоценными камнями и стильные изделия на каждый день. Приглашаем вас познакомиться с миром красоты, элегантности и безупречного качества.",
			image: "/images/news/openstore.jpg",
		},
		{
			title: "Нам исполняется год: праздник для наших клиентов",
			description: "Наш ювелирный магазин отмечает день рождения! В честь этого события мы подготовили специальные предложения, приятные подарки и праздничную атмосферу для всех гостей. Благодарим вас за доверие и будем рады разделить этот особенный день вместе с вами.",
			image: "/images/news/birth.jpg",
		},
		{
			title: "Новогоднее настроение и сияние украшений",
			description: "Новый год уже близко! Самое время выбрать подарки, которые будут радовать и удивлять. В нашем магазине вы найдёте украшения, способные подчеркнуть волшебство праздника и стать незабываемым подарком для ваших близких.",
			image: "/images/news/newyaer.jpg",
		},
	}

	for _, n := range news {
		_, err := DB.Exec(`INSERT OR IGNORE INTO news (title, description, image) VALUES (?, ?, ?)`, n.title, n.description, n.image)
		if err != nil {
			log.Printf("Failed to insert news item: %v", err)
		}
	}

	
	log.Println("Jewelry store data seeded successfully")

	// Создадим администратора по умолчанию если его нет
	var adminID int
	err := DB.QueryRow("SELECT id FROM users WHERE email = ?", "admin@example.com").Scan(&adminID)
	if err == sql.ErrNoRows {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err := DB.Exec("INSERT INTO users (email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", "admin@example.com", string(hashed), "admin", time.Now(), time.Now())
		if err != nil {
			log.Printf("Failed to create admin user: %v", err)
		} else {
			log.Println("Default admin created: admin@example.com / admin123")
		}
	}
}
