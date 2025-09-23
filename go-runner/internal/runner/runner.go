// package runner

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"os/signal"
// 	"path/filepath"
// 	"syscall"
// 	"time"

// 	"github.com/streadway/amqp"
// 	"go-runner/internal/config"
// )

// type Runner struct {
// 	config *config.Config
// }

// type Job struct {
// 	ID        string `json:"id"`
// 	Script    string `json:"script"`
// 	SubDomain string `json:"subdomain"`
// }

// type Response struct {
// 	ID      string `json:"id"`
// 	Status  string `json:"status"`
// 	Message string `json:"message"`
// 	Log     string `json:"log"`
// }

// func New() (*Runner, error) {
// 	cfg, err := config.Load()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Runner{config: cfg}, nil
// }

// func (r *Runner) Start() error {
// 	log.Println("🚀 Starting go-runner...")
// 	log.Printf("📋 Config: Server=%s, Queue=%s, Scripts=%s", 
// 		r.config.Runner.ServerURL, r.config.Runner.MessageQueue, r.config.Runner.ScriptsPath)

// 	// Start consuming jobs
// 	go r.consumeJobs()

// 	// Wait for shutdown signal
// 	r.waitForShutdown()
// 	return nil
// }

// func (r *Runner) consumeJobs() {
// 	for {
// 		if err := r.connectAndConsume(); err != nil {
// 			log.Printf("❌ Consumer error: %v", err)
// 			log.Println("� Retrying in 5 seconds...")
// 			time.Sleep(5 * time.Second)
// 		}
// 	}
// }

// func (r *Runner) connectAndConsume() error {
// 	rabbitURL := fmt.Sprintf("amqp://user:password@%s:5672/", r.config.Runner.MessageQueue)
// 	conn, err := amqp.Dial(rabbitURL)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return fmt.Errorf("failed to open channel: %w", err)
// 	}
// 	defer ch.Close()

// 	queueName := r.config.Runner.ID
// 	if queueName == "" {
// 		queueName = "default-runner"
// 	}

// 	// Declare queue
// 	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to declare queue: %w", err)
// 	}

// 	// Consume messages
// 	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to register consumer: %w", err)
// 	}

// 	log.Printf("👂 Listening on queue: %s", queueName)

// 	for msg := range msgs {
// 		go r.handleJob(msg.Body)
// 	}

// 	return nil
// }

// func (r *Runner) handleJob(body []byte) {
// 	log.Println("📥 Processing job...")

// 	var job Job
// 	if err := json.Unmarshal(body, &job); err != nil {
// 		log.Printf("❌ Failed to unmarshal job: %v", err)
// 		return
// 	}

// 	log.Printf("� Script: %s | Subdomain: %s", job.Script, job.SubDomain)

// 	success, output := r.executeScript(job.Script, job.SubDomain)

// 	response := Response{
// 		ID:      job.ID,
// 		Status:  "done",
// 		Message: fmt.Sprintf("Script '%s' executed", job.Script),
// 		Log:     output,
// 	}

// 	if !success {
// 		response.Status = "error"
// 	}

// 	r.sendResponse(job.ID, response)
// }

// func (r *Runner) executeScript(script, subDomain string) (bool, string) {
// 	log.Printf("▶️ Executing: %s", script)

// 	scriptPath := filepath.Join(r.config.Runner.ScriptsPath, script)
	
// 	// Make script executable
// 	if err := os.Chmod(scriptPath, 0755); err != nil {
// 		log.Printf("⚠️ Failed to make script executable: %v", err)
// 	}

// 	cmd := exec.Command(scriptPath, subDomain)
// 	cmd.Dir = r.config.Runner.ScriptsPath

// 	output, err := cmd.CombinedOutput()
	
// 	if err != nil {
// 		log.Printf("❌ Script failed: %v", err)
// 		return false, string(output)
// 	}

// 	log.Printf("✅ Script completed successfully")
// 	return true, string(output)
// }

// func (r *Runner) sendResponse(jobID string, response Response) {
// 	rabbitURL := fmt.Sprintf("amqp://guest:guest@%s:5672/", r.config.Runner.MessageQueue)
// 	conn, err := amqp.Dial(rabbitURL)
// 	if err != nil {
// 		log.Printf("❌ Failed to connect for response: %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Printf("❌ Failed to open channel for response: %v", err)
// 		return
// 	}
// 	defer ch.Close()

// 	responseQueue := fmt.Sprintf("%s_response", jobID)
	
// 	_, err = ch.QueueDeclare(responseQueue, true, false, false, false, nil)
// 	if err != nil {
// 		log.Printf("❌ Failed to declare response queue: %v", err)
// 		return
// 	}

// 	body, err := json.Marshal(response)
// 	if err != nil {
// 		log.Printf("❌ Failed to marshal response: %v", err)
// 		return
// 	}

// 	err = ch.Publish("", responseQueue, false, false, amqp.Publishing{
// 		ContentType: "application/json",
// 		Body:        body,
// 	})
	
// 	if err != nil {
// 		log.Printf("❌ Failed to publish response: %v", err)
// 		return
// 	}

// 	log.Printf("📨 Response sent to: %s", responseQueue)
// }

// func (r *Runner) waitForShutdown() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

// 	<-c
// 	log.Println("🛑 Shutting down runner...")
// 	log.Println("✅ Runner stopped")
// }







// package runner

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"os/signal"
// 	"path/filepath"
// 	"strings"
// 	"syscall"
// 	"time"

// 	"github.com/streadway/amqp"
// 	"go-runner/internal/config"
// )

// type Runner struct {
// 	config *config.Config
// }

// type Job struct {
// 	ID        string `json:"id"`
// 	Script    string `json:"script"`
// 	SubDomain string `json:"subdomain"`
// }

// type Response struct {
// 	ID      string `json:"id"`
// 	Status  string `json:"status"`
// 	Message string `json:"message"`
// 	Log     string `json:"log"`
// }

// type RegisterRequest struct {
// 	Name string `json:"name"`
// 	IP   string `json:"ip"`
// 	Tags string `json:"tags"`
// }

// type RegisterResponse struct {
// 	Runner struct {
// 		ID string `json:"id"`
// 	} `json:"runner"`
// }

// func New() (*Runner, error) {
// 	cfg, err := config.Load()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Runner{config: cfg}, nil
// }

// func (r *Runner) Start() error {
// 	log.Println("🚀 Starting go-runner...")
// 	log.Printf("📋 Config: Server=%s, Queue=%s, Scripts=%s", 
// 		r.config.Runner.ServerURL, r.config.Runner.MessageQueue, r.config.Runner.ScriptsPath)

// 	// Register with controller first
// 	if err := r.register(); err != nil {
// 		return fmt.Errorf("failed to register runner: %w", err)
// 	}

// 	// Start consuming jobs
// 	go r.consumeJobs()

// 	// Wait for shutdown signal
// 	r.waitForShutdown()
// 	return nil
// }

// func (r *Runner) consumeJobs() {
// 	for {
// 		if err := r.connectAndConsume(); err != nil {
// 			log.Printf("❌ Consumer error: %v", err)
// 			log.Println("� Retrying in 5 seconds...")
// 			time.Sleep(5 * time.Second)
// 		}
// 	}
// }

// func (r *Runner) connectAndConsume() error {
// 	rabbitURL := fmt.Sprintf("amqp://user:password@%s:5672/", r.config.Runner.MessageQueue)
// 	conn, err := amqp.Dial(rabbitURL)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return fmt.Errorf("failed to open channel: %w", err)
// 	}
// 	defer ch.Close()

// 	queueName := r.config.Runner.ID
// 	if queueName == "" {
// 		queueName = "default-runner"
// 	}

// 	// Declare queue
// 	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to declare queue: %w", err)
// 	}

// 	// Consume messages
// 	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to register consumer: %w", err)
// 	}

// 	log.Printf("👂 Listening on queue: %s", queueName)

// 	for msg := range msgs {
// 		go r.handleJob(msg.Body)
// 	}

// 	return nil
// }

// func (r *Runner) handleJob(body []byte) {
// 	log.Println("📥 Processing job...")

// 	var job Job
// 	if err := json.Unmarshal(body, &job); err != nil {
// 		log.Printf("❌ Failed to unmarshal job: %v", err)
// 		return
// 	}

// 	log.Printf("� Script: %s | Subdomain: %s", job.Script, job.SubDomain)

// 	success, output := r.executeScript(job.Script, job.SubDomain)

// 	response := Response{
// 		ID:      job.ID,
// 		Status:  "done",
// 		Message: fmt.Sprintf("Script '%s' executed", job.Script),
// 		Log:     output,
// 	}

// 	if !success {
// 		response.Status = "error"
// 	}

// 	r.sendResponse(job.ID, response)
// }

// func (r *Runner) executeScript(script, subDomain string) (bool, string) {
// 	log.Printf("▶️ Executing: %s", script)

// 	scriptPath := filepath.Join(r.config.Runner.ScriptsPath, script)
	
// 	// Make script executable
// 	if err := os.Chmod(scriptPath, 0755); err != nil {
// 		log.Printf("⚠️ Failed to make script executable: %v", err)
// 	}

// 	cmd := exec.Command(scriptPath, subDomain)
// 	cmd.Dir = r.config.Runner.ScriptsPath

// 	output, err := cmd.CombinedOutput()
	
// 	if err != nil {
// 		log.Printf("❌ Script failed: %v", err)
// 		return false, string(output)
// 	}

// 	log.Printf("✅ Script completed successfully")
// 	return true, string(output)
// }

// func (r *Runner) sendResponse(jobID string, response Response) {
// 	rabbitURL := fmt.Sprintf("amqp://user:password@%s:5672/", r.config.Runner.MessageQueue)
// 	conn, err := amqp.Dial(rabbitURL)
// 	if err != nil {
// 		log.Printf("❌ Failed to connect for response: %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Printf("❌ Failed to open channel for response: %v", err)
// 		return
// 	}
// 	defer ch.Close()

// 	responseQueue := fmt.Sprintf("%s_response", jobID)
	
// 	_, err = ch.QueueDeclare(responseQueue, true, false, false, false, nil)
// 	if err != nil {
// 		log.Printf("❌ Failed to declare response queue: %v", err)
// 		return
// 	}

// 	body, err := json.Marshal(response)
// 	if err != nil {
// 		log.Printf("❌ Failed to marshal response: %v", err)
// 		return
// 	}

// 	err = ch.Publish("", responseQueue, false, false, amqp.Publishing{
// 		ContentType: "application/json",
// 		Body:        body,
// 	})
	
// 	if err != nil {
// 		log.Printf("❌ Failed to publish response: %v", err)
// 		return
// 	}

// 	log.Printf("📨 Response sent to: %s", responseQueue)
// }

// func (r *Runner) waitForShutdown() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

// 	<-c
// 	log.Println("🛑 Shutting down runner...")
// 	log.Println("✅ Runner stopped")
// }

// func (r *Runner) register() error {
// 	// Check if already registered
// 	if r.config.Runner.ID != "default-runner" && r.config.Runner.ID != "" {
// 		log.Println("✅ Runner already registered with ID:", r.config.Runner.ID)
// 		return nil
// 	}

// 	log.Println("🛰 Registering runner with controller...")

// 	// Get local IP for registration
// 	localIP := getLocalIP()
	
// 	payload := RegisterRequest{
// 		Name: r.config.Runner.Name,
// 		IP:   localIP,
// 		Tags: strings.Join(r.config.Runner.Tags, ","),
// 	}

// 	// Try registration multiple times
// 	for i := 0; i < 10; i++ {
// 		if err := r.attemptRegistration(payload); err != nil {
// 			log.Printf("⚠️ Registration attempt %d/10 failed: %v", i+1, err)
// 			time.Sleep(2 * time.Second)
// 			continue
// 		}
// 		return nil
// 	}

// 	return fmt.Errorf("failed to register runner after 10 attempts")
// }

// func (r *Runner) attemptRegistration(payload RegisterRequest) error {
// 	client := &http.Client{Timeout: 5 * time.Second}
	
// 	jsonData, err := json.Marshal(payload)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal payload: %w", err)
// 	}

// 	registerURL := r.config.Runner.ServerURL + "/register"
// 	resp, err := client.Post(registerURL, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return fmt.Errorf("HTTP request failed: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		return fmt.Errorf("server returned error (%d): %s", resp.StatusCode, string(body))
// 	}

// 	var registerResp RegisterResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
// 		return fmt.Errorf("failed to decode response: %w", err)
// 	}

// 	// Update runner ID
// 	r.config.Runner.ID = registerResp.Runner.ID
	
// 	// Save ID to .env file
// 	if err := r.saveRunnerID(registerResp.Runner.ID); err != nil {
// 		log.Printf("⚠️ Failed to save RUNNER_ID to .env: %v", err)
// 	}

// 	log.Printf("✅ Runner registered successfully! ID: %s", registerResp.Runner.ID)
// 	return nil
// }

// func (r *Runner) saveRunnerID(runnerID string) error {
// 	envContent, err := ioutil.ReadFile(".env")
// 	if err != nil && !os.IsNotExist(err) {
// 		return err
// 	}

// 	lines := strings.Split(string(envContent), "\n")
// 	runnerIDFound := false

// 	for i, line := range lines {
// 		if strings.HasPrefix(line, "RUNNER_ID=") {
// 			lines[i] = fmt.Sprintf("RUNNER_ID=%s", runnerID)
// 			runnerIDFound = true
// 			break
// 		}
// 	}

// 	if !runnerIDFound {
// 		lines = append(lines, fmt.Sprintf("RUNNER_ID=%s", runnerID))
// 	}

// 	newContent := strings.Join(lines, "\n")
// 	return ioutil.WriteFile(".env", []byte(newContent), 0644)
// }

// // getLocalIP helper function
// func getLocalIP() string {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return "127.0.0.1"
// 	}

// 	for _, addr := range addrs {
// 		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if ipnet.IP.To4() != nil {
// 				return ipnet.IP.String()
// 			}
// 		}
// 	}

// 	return "127.0.0.1"
// }













package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/streadway/amqp"
	"go-runner/internal/config"
)

type Runner struct {
	config *config.Config
}

type Job struct {
	ID        string `json:"id"`
	Script    string `json:"script"`
	SubDomain string `json:"subdomain"`
}

type Response struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Log     string `json:"log"`
}

type RegisterRequest struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Tags string `json:"tags"`
}

type RegisterResponse struct {
	Runner struct {
		ID string `json:"id"`
	} `json:"runner"`
}

func New() (*Runner, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	return &Runner{config: cfg}, nil
}

func (r *Runner) Start() error {
	log.Println("Starting go-runner...")
	log.Printf("Config: Server=%s, Queue=%s, Scripts=%s", 
		r.config.Runner.ServerURL, r.config.Runner.MessageQueue, r.config.Runner.ScriptsPath)

	// Register with controller first
	if err := r.register(); err != nil {
		return fmt.Errorf("failed to register runner: %w", err)
	}

	// Start consuming jobs
	go r.consumeJobs()

	// Wait for shutdown signal
	r.waitForShutdown()
	return nil
}

func (r *Runner) consumeJobs() {
	for {
		if err := r.connectAndConsume(); err != nil {
			log.Printf("Consumer error: %v", err)
			log.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}

func (r *Runner) connectAndConsume() error {
	rabbitURL := fmt.Sprintf("amqp://user:password@%s:5672/", r.config.Runner.MessageQueue)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	queueName := r.config.Runner.ID
	if queueName == "" {
		queueName = "default-runner"
	}

	// Declare queue
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Consume messages
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Listening on queue: %s", queueName)

	for msg := range msgs {
		go r.handleJob(msg.Body)
	}

	return nil
}

func (r *Runner) handleJob(body []byte) {
	log.Println("Processing job...")

	var job Job
	if err := json.Unmarshal(body, &job); err != nil {
		log.Printf("Failed to unmarshal job: %v", err)
		return
	}

	log.Printf("Script: %s | Subdomain: %s", job.Script, job.SubDomain)

	success, output := r.executeScript(job.Script, job.SubDomain)

	response := Response{
		ID:      job.ID,
		Status:  "done",
		Message: fmt.Sprintf("Script '%s' executed", job.Script),
		Log:     output,
	}

	if !success {
		response.Status = "error"
	}

	r.sendResponse(job.ID, response)
}

func (r *Runner) executeScript(script, subDomain string) (bool, string) {
	log.Printf("Executing: %s", script)

	scriptPath := filepath.Join(r.config.Runner.ScriptsPath, script)
	
	// Make script executable
	// if err := os.Chmod(scriptPath, 0755); err != nil {
	// 	log.Printf("Failed to make script executable: %v", err)
	// }

	cmd := exec.Command(scriptPath, subDomain)
	cmd.Dir = r.config.Runner.ScriptsPath

	output, err := cmd.CombinedOutput()
	
	if err != nil {
		log.Printf("Script failed: %v", err)
		return false, string(output)
	}

	log.Printf("Script completed successfully")
	return true, string(output)
}

func (r *Runner) sendResponse(jobID string, response Response) {
	rabbitURL := fmt.Sprintf("amqp://user:password@%s:5672/", r.config.Runner.MessageQueue)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Printf("Failed to connect for response: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel for response: %v", err)
		return
	}
	defer ch.Close()

	responseQueue := fmt.Sprintf("%s_response", jobID)
	
	_, err = ch.QueueDeclare(responseQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare response queue: %v", err)
		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	err = ch.Publish("", responseQueue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	
	if err != nil {
		log.Printf("Failed to publish response: %v", err)
		return
	}

	log.Printf("Response sent to: %s", responseQueue)
}

func (r *Runner) waitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Shutting down runner...")
	log.Println("Runner stopped")
}

func (r *Runner) register() error {
	// Check if already registered (KEY_ID exists)
	if r.config.Runner.ID != "" {
		log.Println("Runner already registered with KEY_ID:", r.config.Runner.ID)
		return nil
	}

	log.Println("🛰 Registering runner with controller...")

	// Get local IP for registration
	localIP := getLocalIP()
	
	payload := RegisterRequest{
		Name: r.config.Runner.Name,
		IP:   localIP,
		Tags: strings.Join(r.config.Runner.Tags, ","),
	}

	// Try registration multiple times
	for i := 0; i < 10; i++ {
		if err := r.attemptRegistration(payload); err != nil {
			log.Printf("Registration attempt %d/10 failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to register runner after 10 attempts")
}

func (r *Runner) attemptRegistration(payload RegisterRequest) error {
	client := &http.Client{Timeout: 5 * time.Second}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	registerURL := r.config.Runner.ServerURL + "/register"
	resp, err := client.Post(registerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("server returned error (%d): %s", resp.StatusCode, string(body))
	}

	var registerResp RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Update runner ID
	r.config.Runner.ID = registerResp.Runner.ID
	
	// Save KEY_ID to .env file (for compatibility with Python runner)
	if err := r.saveKeyID(registerResp.Runner.ID); err != nil {
		log.Printf("Failed to save KEY_ID to .env: %v", err)
	}

	log.Printf("Runner registered successfully! ID: %s", registerResp.Runner.ID)
	return nil
}

func (r *Runner) saveKeyID(keyID string) error {
	envContent, err := ioutil.ReadFile(".env")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(string(envContent), "\n")
	keyIDFound := false

	for i, line := range lines {
		if strings.HasPrefix(line, "KEY_ID=") {
			lines[i] = fmt.Sprintf("KEY_ID=%s", keyID)
			keyIDFound = true
			break
		}
	}

	if !keyIDFound {
		lines = append(lines, fmt.Sprintf("KEY_ID=%s", keyID))
	}

	newContent := strings.Join(lines, "\n")
	return ioutil.WriteFile(".env", []byte(newContent), 0644)
}

// getLocalIP helper function
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}

