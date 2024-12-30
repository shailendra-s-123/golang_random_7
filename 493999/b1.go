package main  
import (  
    "fmt"
    "net/url"
    "reflect"
    "strconv"
    "strings"
)

type Params struct {
    data map[string][]interface{}
}

func NewParams() *Params {
    return &Params{
        data: make(map[string][]interface{}),
    }
}
func (p *Params) Set(key string, value interface{}) {
    if _, ok := p.data[key]; !ok {
        p.data[key] = []interface{}{value}
    } else {
        p.data[key] = append(p.data[key], value)
    }
}
func (p *Params) Get(key string) []interface{} {
    return p.data[key]
}
func (p *Params) String(key string) string {
    v := p.Get(key)
    if len(v) == 0 {
        return ""
    }
    return fmt.Sprintf("%v", v[0])
}

func (p *Params) Int(key string) (int, error) {
    s := p.String(key)
    if s == "" {
        return 0, nil
    }
    return strconv.Atoi(s)
}

func (p *Params) Bool(key string) (bool, error) {
    s := p.String(key)
    if s == "" {
        return false, nil
    }
    return strconv.ParseBool(s)
}

func (p *Params) Float64(key string) (float64, error) {  
    s := p.String(key)
    if s == "" {
        return 0, nil
    }
    return strconv.ParseFloat(s, 64)
}
func (p *Params) ParseQuery(query string) error {
    u, err := url.ParseQuery(query)
    if err != nil {
        return err
    }
    for key, values := range u {
        for _, value := range values {
            v := reflect.ValueOf(value)
            switch v.Kind() {
            case reflect.String:
                p.Set(key, value)
            case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                i, _ := strconv.Atoi(value)
                p.Set(key, i)
            case reflect.Float32, reflect.Float64:
                f, _ := strconv.ParseFloat(value, 64)
                p.Set(key, f)
            case reflect.Bool:
                b, _ := strconv.ParseBool(value)
                p.Set(key, b)
            default:
                p.Set(key, value)
            }
        }
    }
    return nil
}
type DataInsight struct {  
    Param   string
    Values []interface{}
}

func (p *Params) GenerateInsights() []DataInsight {
    insights := make([]DataInsight, 0, len(p.data))
    for key, values := range p.data {
        insights = append(insights, DataInsight{Param: key, Values: values})
    }
    return insights
}
func main() {  
    params := NewParams()
    if err := params.ParseQuery("name=Alice&age=25&is_active=true&tags=sport&tags=music&gender=female"); err != nil {
        fmt.Println("Error parsing query:", err)
        return
    }
    fmt.Println("All Parameters:")
    fmt.Println(params.GenerateInsights())

    name := params.String("name")
    age, _ := params.Int("age")
    isActive, _ := params.Bool("is_active")

    fmt.Println("\nParsed Data:")
    fmt.Println("Name:", name)
    fmt.Println("Age:", age)
    fmt.Println("IsActive:", isActive)

    fmt.Println("\nFiltered Tags:")
    tags := params.Get("tags")
    for _, tag := range tags {
        if strings.HasPrefix(tag.(string), "sp") {
            fmt.Println(tag)
        }
    }
}