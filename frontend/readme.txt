# 创建前端目录结构
New-Item -ItemType Directory -Path "frontend/src/assets"
New-Item -ItemType Directory -Path "frontend/src/components"
New-Item -ItemType Directory -Path "frontend/src/router"
New-Item -ItemType Directory -Path "frontend/src/views"
New-Item -ItemType Directory -Path "frontend/src/utils"
New-Item -ItemType Directory -Path "frontend/src/stores"
New-Item -ItemType Directory -Path "frontend/public"

# 创建后端目录结构
New-Item -ItemType Directory -Path "backend/internal/app/middleware"
New-Item -ItemType Directory -Path "backend/internal/app/router"
New-Item -ItemType Directory -Path "backend/internal/app/service"
New-Item -ItemType Directory -Path "backend/internal/app/repository"
New-Item -ItemType Directory -Path "backend/internal/app/model"
New-Item -ItemType Directory -Path "backend/internal/app/config"
New-Item -ItemType Directory -Path "backend/internal/pkg/utils"
New-Item -ItemType Directory -Path "backend/internal/pkg/constants"
New-Item -ItemType Directory -Path "backend/api"
New-Item -ItemType Directory -Path "backend/scripts"
New-Item -ItemType Directory -Path "backend/tests/unit"
New-Item -ItemType Directory -Path "backend/tests/integration"







#启动脚本
cd e:\workplace\newItem\frontend
npm install
npm run dev