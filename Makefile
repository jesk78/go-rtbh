GO_RTBH = go-rtbh
GO_RTBHAPP = go-rtbhapp

BUILD_DIR = ./build
COMMAND_DIR = ./commands
RESOURCES_DIR = ./api/resources

all: clean $(GO_RTBH) $(GO_RTBHAPP) ${GO_BIRDAPP}

$(GO_RTBH):
	[ ! -d $(BUILD_DIR) ] && mkdir -p $(BUILD_DIR) || true
	go build -v -o $(BUILD_DIR)/$(GO_RTBH) $(COMMAND_DIR)/$(GO_RTBH).go

$(GO_RTBHAPP):
	[ ! -d $(BUILD_DIR) ] && mkdir -p $(BUILD_DIR) || true
	go build -v -o $(BUILD_DIR)/$(GO_RTBHAPP) $(COMMAND_DIR)/$(GO_RTBHAPP).go

$(RTBH_APP):
	cd api/resources/app ; npm install

clean:
	rm -vf $(BUILD_DIR)/$(GO_RTBH)
	rm -vf $(BUILD_DIR)/$(GO_RTBHAPP)
	rm -vrf $(RESOURCES_DIR)/app/node_modules
	rm -vrf $(RESOURCES_DIR)/app/typings
	rm -vf $(RESOURCES_DIR)/app/*.js
	rm -vf $(RESOURCES_DIR)/app/*.map
