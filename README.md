# go-openai-tools

This is a helper library intended to be used with [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai).

## Example

```go
type GetWeatherRequest struct{
    Location string `json:"location" required:"true" description:"City and country e.g. Bogotá, Colombia"`
}

type GetWeatherResponse struct{
    Weather string `json:"weather"`
}

func GetWeather(_ GetWeatherRequest) (GetWeatherResponse, error) {
    return GetWeatherResponse{Weather: "sunny"}, nil
}
```
