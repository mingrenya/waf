#!/bin/bash

find . -name "*.go" -type f -exec sed -i \
-e 's|"github.com/your-org/ai-waf/backend|"coraza-waf/backend|g' \
-e 's|"ai-waf/backend|"coraza-waf/backend|g' \
-e 's|"coraza-waf/backend/backend|"coraza-waf/backend|g' \
{} +
