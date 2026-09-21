package providers

func stringValue[T ~string](v *T) string {
    if v == nil {
        return ""
    }
    return string(*v)
}

type missingEnvironmentError struct {
    Name string
}

func (e *missingEnvironmentError) Error() string {
    return "required environment variable is not set: " + e.Name
}
