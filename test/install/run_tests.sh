#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Set working directory to project root
cd "$(dirname "$0")/../.."
PROJECT_ROOT=$(pwd)

# Default configurations
DISTROS=("ubuntu" "fedora" "centos" "arch" "opensuse")
TEST_ALL=true
KEEP_CONTAINERS=false
VERBOSE=false

# Display help message
function show_help {
    echo "Usage: $0 [OPTIONS]"
    echo "Run installation tests in containerized environments"
    echo ""
    echo "Options:"
    echo "  -d, --distro DISTRO     Test specific distribution (ubuntu, fedora, centos, arch, opensuse)"
    echo "  -k, --keep              Keep containers after tests"
    echo "  -v, --verbose           Show verbose output"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                      Test on all distributions"
    echo "  $0 -d ubuntu            Test only on Ubuntu"
    echo "  $0 -d fedora -k         Test on Fedora and keep the container"
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--distro)
            DISTRO="$2"
            TEST_ALL=false
            shift 2
            ;;
        -k|--keep)
            KEEP_CONTAINERS=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Verify the selected distro is valid
if [ "$TEST_ALL" = false ]; then
    VALID_DISTRO=false
    for d in "${DISTROS[@]}"; do
        if [ "$d" = "$DISTRO" ]; then
            VALID_DISTRO=true
            break
        fi
    done
    
    if [ "$VALID_DISTRO" = false ]; then
        echo -e "${RED}Error: Invalid distribution '$DISTRO'${NC}"
        echo "Valid options are: ${DISTROS[*]}"
        exit 1
    fi
    
    # Set the distros array to only include the selected distro
    DISTROS=("$DISTRO")
fi

# Function to run test for a distribution
function run_test {
    local distro=$1
    local container_name="mbvpn-test-$distro"
    local dockerfile="test/install/dockerfiles/Dockerfile.$distro"
    
    echo -e "\n${BLUE}=======================================${NC}"
    echo -e "${BLUE}Testing installation on $distro${NC}"
    echo -e "${BLUE}=======================================${NC}"
    
    # Build the Docker image
    echo -e "\n${YELLOW}Building Docker image for $distro...${NC}"
    docker build -t "mbvpn-test-$distro" -f "$dockerfile" .
    
    # Run the container
    echo -e "\n${YELLOW}Running test in $distro container...${NC}"
    docker run --rm -d --name "$container_name" \
        -v "$PROJECT_ROOT:/mbvpn" \
        "mbvpn-test-$distro" tail -f /dev/null
    
    # Run the test script inside the container
    if [ "$VERBOSE" = true ]; then
        docker exec "$container_name" bash /mbvpn/test/install/test_install.sh
    else
        docker exec "$container_name" bash /mbvpn/test/install/test_install.sh > /dev/null
    fi
    
    # Check the result
    if [ $? -eq 0 ]; then
        echo -e "\n${GREEN}✅ Installation test on $distro passed!${NC}"
        TEST_RESULT=0
    else
        echo -e "\n${RED}❌ Installation test on $distro failed!${NC}"
        TEST_RESULT=1
    fi
    
    # Clean up container if not keeping
    if [ "$KEEP_CONTAINERS" = false ]; then
        docker stop "$container_name" > /dev/null
    else
        echo -e "\nContainer '$container_name' is still running for inspection."
        echo "To stop it later, run: docker stop $container_name"
    fi
    
    return $TEST_RESULT
}

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed or not in PATH${NC}"
    echo "Please install Docker before running this test suite."
    exit 1
fi

# Run tests for all specified distributions
FAILED_DISTROS=()
PASSED_DISTROS=()

for distro in "${DISTROS[@]}"; do
    if run_test "$distro"; then
        PASSED_DISTROS+=("$distro")
    else
        FAILED_DISTROS+=("$distro")
    fi
done

# Print summary
echo -e "\n${BLUE}=======================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}=======================================${NC}"

if [ ${#PASSED_DISTROS[@]} -gt 0 ]; then
    echo -e "\n${GREEN}Passed (${#PASSED_DISTROS[@]})${NC}"
    for distro in "${PASSED_DISTROS[@]}"; do
        echo -e "  ${GREEN}✅ $distro${NC}"
    done
fi

if [ ${#FAILED_DISTROS[@]} -gt 0 ]; then
    echo -e "\n${RED}Failed (${#FAILED_DISTROS[@]})${NC}"
    for distro in "${FAILED_DISTROS[@]}"; do
        echo -e "  ${RED}❌ $distro${NC}"
    done
    exit 1
fi

echo -e "\n${GREEN}All tests passed successfully!${NC}"
exit 0