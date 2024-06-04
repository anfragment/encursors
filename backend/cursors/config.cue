AllowLocalhost: bool | *false
MinEventTimeoutMs: 1500

if #Meta.Environment.Type == "development" {
  AllowLocalhost: true
}