# Go Patterns and Idioms

A collection of common Go patterns and best practices used in production code.

## 1. Sentinel Error Pattern

Use sentinel errors for errors that callers might want to check programmatically.

```go
import "errors"

var (
    // ErrNotFound is returned when a resource is not found.
    ErrNotFound = errors.New("resource not found")
    
    // ErrUnauthorized is returned when access is denied.
    ErrUnauthorized = errors.New("unauthorized")
)

func FindUser(id string) (*User, error) {
    user, err := db.Get(id)
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    return user, err
}

// Callers can check:
if errors.Is(err, ErrNotFound) {
    // handle not found
}
```

**When to use:**
- Errors that callers might want to check programmatically
- Errors representing specific conditions
- Errors that should be documented as part of the package API

---

## 2. Options Pattern (Functional Options)

For flexible API configuration without many constructor variants.

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
    logger  Logger
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) { s.host = host }
}

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) { s.timeout = timeout }
}

func WithLogger(logger Logger) Option {
    return func(s *Server) { s.logger = logger }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        host:    "localhost",
        port:    8080,
        timeout: 30 * time.Second,
        logger:  defaultLogger,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage:
server := NewServer(
    WithHost("0.0.0.0"),
    WithPort(9000),
    WithTimeout(60*time.Second),
)
```

**Benefits:**
- Backward compatible (can add new options)
- Self-documenting (option names are clear)
- Flexible (any combination of options)

---

## 3. Interface Segregation Pattern

Prefer small, focused interfaces over large ones.

```go
// ✅ GOOD - Small, focused interfaces
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

type Closer interface {
    Close() error
}

// Compose when needed
type ReadWriter interface {
    Reader
    Writer
}

// ❌ BAD - Large interface
type ReadWriteCloser interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
    Flush() error
    Seek(int64, int) (int64, error)
}

// Real-world example from io package:
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}
```

**Benefits:**
- Easier to implement
- More flexible
- Better testability

---

## 4. Context Pattern

For cancellation, timeouts, and request-scoped values.

```go
import "context"

func DoWork(ctx context.Context) error {
    // Check if cancelled
    if err := ctx.Err(); err != nil {
        return err
    }
    
    // Pass context to operations
    return someOperation(ctx)
}

// Usage with timeout:
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := DoWork(ctx); err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // handle timeout
    }
}

// Usage with cancellation:
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    time.Sleep(2 * time.Second)
    cancel() // Cancel after 2 seconds
}()

// Usage with values:
type key string
const userIDKey key = "userID"

ctx := context.WithValue(context.Background(), userIDKey, "123")
userID := ctx.Value(userIDKey).(string)
```

**Best practices:**
- Always pass context as first parameter
- Check `ctx.Err()` before long operations
- Use context for cancellation, not for passing optional parameters

---

## 5. Builder Pattern

For complex object construction with method chaining.

```go
type Query struct {
    table  string
    where  []string
    orderBy string
    limit  int
}

type QueryBuilder struct {
    query Query
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{}
}

func (qb *QueryBuilder) Table(name string) *QueryBuilder {
    qb.query.table = name
    return qb
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
    qb.query.where = append(qb.query.where, condition)
    return qb
}

func (qb *QueryBuilder) OrderBy(field string) *QueryBuilder {
    qb.query.orderBy = field
    return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
    qb.query.limit = n
    return qb
}

func (qb *QueryBuilder) Build() Query {
    return qb.query
}

// Usage:
query := NewQueryBuilder().
    Table("users").
    Where("age > 18").
    Where("active = true").
    OrderBy("created_at DESC").
    Limit(10).
    Build()
```

---

## 6. Dependency Injection Pattern

Constructor-based dependency injection for testability.

```go
type Repository interface {
    Find(id string) (*User, error)
    Save(user *User) error
}

type Logger interface {
    Info(msg string)
    Error(msg string)
}

type Service struct {
    repo   Repository
    logger Logger
}

// Constructor with dependencies
func NewService(repo Repository, logger Logger) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
    }
}

func (s *Service) GetUser(id string) (*User, error) {
    s.logger.Info("fetching user")
    return s.repo.Find(id)
}

// Usage:
repo := NewPostgresRepo()
logger := NewLogger()
service := NewService(repo, logger)

// Testing:
mockRepo := &MockRepository{}
mockLogger := &MockLogger{}
testService := NewService(mockRepo, mockLogger)
```

**Benefits:**
- Testable (can inject mocks)
- Flexible (can swap implementations)
- Clear dependencies

---

## 7. Worker Pool Pattern

For concurrent processing with controlled parallelism.

```go
type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID int
    Output string
    Err   error
}

func ProcessJobs(jobs <-chan Job, results chan<- Result, workers int) {
    var wg sync.WaitGroup
    
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                output, err := processJob(job)
                results <- Result{
                    JobID:  job.ID,
                    Output: output,
                    Err:    err,
                }
            }
        }(i)
    }
    
    wg.Wait()
    close(results)
}

// Usage:
jobs := make(chan Job, 100)
results := make(chan Result, 100)

// Send jobs
go func() {
    for i := 0; i < 1000; i++ {
        jobs <- Job{ID: i, Data: fmt.Sprintf("data-%d", i)}
    }
    close(jobs)
}()

// Process with 10 workers
go ProcessJobs(jobs, results, 10)

// Collect results
for result := range results {
    if result.Err != nil {
        log.Printf("Job %d failed: %v", result.JobID, result.Err)
    } else {
        log.Printf("Job %d: %s", result.JobID, result.Output)
    }
}
```

---

## 8. Error Wrapping Pattern

Using `%w` verb for error chains (Go 1.13+).

```go
func DoSomething() error {
    if err := step1(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    if err := step2(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}

// Check with errors.Is():
if errors.Is(err, ErrNotFound) {
    // handle not found
}

// Unwrap with errors.Unwrap():
originalErr := errors.Unwrap(err)

// Check error type with errors.As():
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // handle PathError
}
```

**Best practices:**
- Use `%w` to wrap errors, not `%v`
- Add context at each level
- Use `errors.Is()` for sentinel errors
- Use `errors.As()` for type assertions

---

## 9. Guard Clauses Pattern

Early returns for cleaner error handling.

```go
// ✅ GOOD - Guard clauses
func Process(data Data) error {
    if err := validate(data); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    if err := transform(data); err != nil {
        return fmt.Errorf("transformation failed: %w", err)
    }
    if err := save(data); err != nil {
        return fmt.Errorf("save failed: %w", err)
    }
    return nil
}

// ❌ BAD - Nested ifs
func Process(data Data) error {
    if err := validate(data); err == nil {
        if err := transform(data); err == nil {
            if err := save(data); err == nil {
                return nil
            } else {
                return err
            }
        } else {
            return err
        }
    } else {
        return err
    }
}
```

**Benefits:**
- Reduces nesting
- Easier to read
- Clear error paths

---

## 10. Interface Assertion Pattern

Type assertions with ok checks.

```go
type Stringer interface {
    String() string
}

func Print(v interface{}) {
    if s, ok := v.(Stringer); ok {
        fmt.Println(s.String())
    } else {
        fmt.Printf("%v\n", v)
    }
}

// Type switch
func Process(v interface{}) {
    switch val := v.(type) {
    case string:
        fmt.Println("string:", val)
    case int:
        fmt.Println("int:", val)
    case bool:
        fmt.Println("bool:", val)
    default:
        fmt.Println("unknown type")
    }
}

// Interface assertion
func AssertWriter(w io.Writer) error {
    if _, ok := w.(io.Closer); ok {
        return fmt.Errorf("writer is also a closer")
    }
    return nil
}
```

---

## 11. Once Pattern (sync.Once)

For one-time initialization (thread-safe singleton).

```go
var (
    instance *Service
    once     sync.Once
)

func GetInstance() *Service {
    once.Do(func() {
        instance = &Service{
            // initialization
        }
    })
    return instance
}

// With initialization function:
var (
    instance *Service
    once     sync.Once
)

func GetInstance() *Service {
    once.Do(initInstance)
    return instance
}

func initInstance() {
    instance = &Service{
        // initialization
    }
}
```

**Use cases:**
- Lazy initialization
- Singleton pattern
- One-time setup

---

## 12. Defer Pattern

For cleanup and resource management.

```go
// File handling
func ProcessFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close() // Always closes, even on panic
    
    // Process file...
    return nil
}

// Multiple defers (LIFO - Last In First Out)
func Example() {
    defer fmt.Println("third")
    defer fmt.Println("second")
    defer fmt.Println("first")
    // Output: first, second, third
}

// Defer with error handling
func DoSomething() (err error) {
    resource, err := Acquire()
    if err != nil {
        return err
    }
    defer func() {
        if closeErr := resource.Close(); closeErr != nil {
            err = closeErr // Modify return value
        }
    }()
    
    // Use resource...
    return nil
}

// Defer for unlocking
func (m *Mutex) LockAndDo() {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    // Critical section
}
```

**Best practices:**
- Always use defer for cleanup
- Defer executes even on panic
- Defers execute in LIFO order

---

## 13. Channel Patterns

### Generator Pattern

```go
func GenerateNumbers(max int) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < max; i++ {
            ch <- i
        }
    }()
    return ch
}

// Usage:
for num := range GenerateNumbers(10) {
    fmt.Println(num)
}
```

### Fan-out/Fan-in Pattern

```go
// Fan-out: Distribute work to multiple workers
func FanOut(input <-chan Job, workers int) []<-chan Job {
    outputs := make([]<-chan Job, workers)
    for i := 0; i < workers; i++ {
        output := make(chan Job)
        outputs[i] = output
        go func() {
            defer close(output)
            for job := range input {
                output <- job
            }
        }()
    }
    return outputs
}

// Fan-in: Combine multiple channels into one
func FanIn(inputs ...<-chan Result) <-chan Result {
    output := make(chan Result)
    var wg sync.WaitGroup
    
    for _, input := range inputs {
        wg.Add(1)
        go func(ch <-chan Result) {
            defer wg.Done()
            for result := range ch {
                output <- result
            }
        }(input)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}
```

---

## 14. Middleware Pattern

For HTTP handlers and cross-cutting concerns.

```go
type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Chain middleware
func Chain(middlewares ...Middleware) Middleware {
    return func(next http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next
    }
}

// Usage:
handler := Chain(
    LoggingMiddleware,
    AuthMiddleware,
)(myHandler)
```

---

## 15. Strategy Pattern

Using interfaces for different algorithms.

```go
type PaymentStrategy interface {
    Pay(amount float64) error
}

type CreditCardPayment struct{}
func (c *CreditCardPayment) Pay(amount float64) error {
    // Process credit card
    return nil
}

type PayPalPayment struct{}
func (p *PayPalPayment) Pay(amount float64) error {
    // Process PayPal
    return nil
}

type PaymentProcessor struct {
    strategy PaymentStrategy
}

func (p *PaymentProcessor) SetStrategy(s PaymentStrategy) {
    p.strategy = s
}

func (p *PaymentProcessor) ProcessPayment(amount float64) error {
    return p.strategy.Pay(amount)
}

// Usage:
processor := &PaymentProcessor{}
processor.SetStrategy(&CreditCardPayment{})
processor.ProcessPayment(100.0)
```

---

## Summary

These patterns help write:
- **Testable** code (DI, interfaces)
- **Maintainable** code (guard clauses, error wrapping)
- **Concurrent** code (worker pools, channels)
- **Flexible** code (options pattern, strategy pattern)
- **Safe** code (defer, sync.Once)

Choose patterns based on your specific needs and context.
