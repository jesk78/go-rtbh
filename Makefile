TARGET = go-rtbh

BUILD_DIR = ./build
CMD_DIR = ./cmd
RESOURCES_DIR = ./api/resources

all: $(BUILD_DIR)/$(TARGET)

$(BUILD_DIR):
	mkdir -p "$(BUILD_DIR)"

$(BUILD_DIR)/$(TARGET):
	go build -v -o $(BUILD_DIR)/$(TARGET) $(CMD_DIR)/$(TARGET)/$(TARGET).go

clean:
	rm -rf $(BUILD_DIR) || true
