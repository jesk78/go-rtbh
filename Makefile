GO_RTBH = go-rtbh
GO_RTBHAPI = go-rtbhapi
RTBH_APP = rtbh-app

BUILD_DIR = ./build
COMMAND_DIR = ./commands
RESOURCES_DIR = ./api/resources

all: clean $(GO_RTBH) $(GO_RTBHAPI)

$(GO_RTBH):
	[ ! -d $(BUILD_DIR) ] && mkdir -p $(BUILD_DIR) || true
	go build -v -o $(BUILD_DIR)/$(GO_RTBH) $(COMMAND_DIR)/$(GO_RTBH)/$(GO_RTBH).go

$(GO_RTBHAPI):
	[ ! -d $(BUILD_DIR) ] && mkdir -p $(BUILD_DIR) || true
	go build -v -o $(BUILD_DIR)/$(GO_RTBHAPI) $(COMMAND_DIR)/$(GO_RTBHAPI)/$(GO_RTBHAPI).go

$(RTBH_APP):
	cd api/resources/app ; npm install

clean:
	rm -vf $(BUILD_DIR)/$(GO_RTBH)
	rm -vf $(BUILD_DIR)/$(GO_RTBHAPI)
	rm -vrf $(RESOURCES_DIR)/app/node_modules
	rm -vrf $(RESOURCES_DIR)/app/typings
	rm -vf $(RESOURCES_DIR)/app/*.js
	rm -vf $(RESOURCES_DIR)/app/*.map
