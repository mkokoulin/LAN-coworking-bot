package services

import (
	// "context"
	// "encoding/json"
	// "fmt"
	"net/http"
	// "net/url"
	// "sort"
	// "strconv"
	"strings"

	// "github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

type HaysellBarService struct {
    client  *http.Client
    baseURL string
    apiKey  string
}

type haysellCategory struct {
    Sort       int    `json:"sort" bson:"sort"`
    Deleted    bool   `json:"deleted" bson:"deleted"`
    Id         int    `json:"id" bson:"id"`
    NameRu     string `json:"nameRu" bson:"name_ru"`
    Name       string `json:"name" bson:"name"`
    Active     int    `json:"active" bson:"active"`
    ParentId   int    `json:"parentId" bson:"parent_id"`
    Level      int    `json:"level" bson:"level"`
}

type haysellProduct struct {
    ID          int    `json:"id"`
    ImagePath   string `json:"imagePath"`
    ItemName    string `json:"itemName"`
    ShortName   string `json:"shortName"`
    Price       int    `json:"price"`
    Category    string `json:"category"`
    Description string `json:"description"`
    Balance     int    `json:"balance"`
}

func NewHaysellBarService(client *http.Client, baseURL, apiKey string) *HaysellBarService {
    return &HaysellBarService{
        client:  client,
        baseURL: strings.TrimRight(baseURL, "/"),
        apiKey:  apiKey,
    }
}

// func (s *HaysellBarService) ListCategories(ctx context.Context) ([]types.BarCategory, error) {
//     req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/api/categories?parentId=1240", nil)
//     if err != nil {
//         return nil, fmt.Errorf("new redquest: %w", err)
//     }

//     // if s.apiKey != "" {
//     //     // если нужен другой заголовок — поменяешь тут
//     //     // req.Header.Set("Authorization", "Bearer "+s.apiKey)
//     // }

//     resp, err := s.client.Do(req)
//     if err != nil {
//         return nil, fmt.Errorf("do request: %w", err)
//     }
//     defer resp.Body.Close()

//     if resp.StatusCode != http.StatusOK {
//         return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
//     }

//     var raw []haysellCategory
//     if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
//         return nil, fmt.Errorf("decode categories: %w", err)
//     }

//     res := make([]types.BarCategory, 0, len(raw))
//     for _, c := range raw {
//         if c.Deleted || c.Active == 0 {
//             continue
//         }
//         id := strconv.Itoa(c.Id)
//         name := c.NameRu
//         if name == "" {
//             name = c.Name
//         }
//         res = append(res, types.BarCategory{
//             ID:   id,
//             Name: name,
//         })
//     }

//     // сортировка по имени
//     sort.Slice(res, func(i, j int) bool {
//         return res[i].Name < res[j].Name
//     })

//     return res, nil
// }

// func (s *HaysellBarService) ListProducts(ctx context.Context, ids []int) ([]types.BarProduct, error) {
//     if len(ids) == 0 {
//         return []types.BarProduct{}, nil
//     }

//     u, err := url.Parse(s.baseURL + "/api/products")
//     if err != nil {
//         return nil, fmt.Errorf("parse base url: %w", err)
//     }

//     q := u.Query()
//     for _, id := range ids {
//         q.Add("id", strconv.Itoa(id))
//     }
//     u.RawQuery = q.Encode()

//     req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
//     if err != nil {
//         return nil, fmt.Errorf("new request: %w", err)
//     }

//     if s.apiKey != "" {
//         // req.Header.Set("Authorization", "Bearer "+s.apiKey)
//     }

//     resp, err := s.client.Do(req)
//     if err != nil {
//         return nil, fmt.Errorf("do request: %w", err)
//     }
//     defer resp.Body.Close()

//     if resp.StatusCode != http.StatusOK {
//         return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
//     }

//     var raw map[string]haysellProduct
//     if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
//         return nil, fmt.Errorf("decode products: %w", err)
//     }

//     res := make([]types.BarProduct, 0, len(raw))
//     for _, p := range raw {
//         res = append(res, types.BarProduct{
//             ID:          p.ID,
//             CategoryID:  p.Category,
//             Name:        p.ItemName,
//             ShortName:   p.ShortName,
//             PriceAMD:    p.Price,
//             ImageURL:    p.ImagePath,
//             Description: p.Description,
//             Balance:     p.Balance,
//         })
//     }

//     sort.Slice(res, func(i, j int) bool {
//         return res[i].Name < res[j].Name
//     })

//     return res, nil
// }
