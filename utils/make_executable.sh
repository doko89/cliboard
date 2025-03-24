#!/bin/bash

# Navigate to project root
cd "$(dirname "$0")/.."

# Make scripts executable
chmod +x utils/fix_dependencies.sh
chmod +x utils/build.sh
chmod +x utils/make_executable.sh
chmod +x utils/fix_go_mod.sh
chmod +x utils/fix_all_dependencies.sh
chmod +x utils/fix_go_sum.sh
chmod +x utils/force_fix_deps.sh

echo "Scripts are now executable."
