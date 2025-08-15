@/home/ubuntu/workspace/vocabulary/vibe-coding-starter-ui-antd/CLAUDE.md
@/home/ubuntu/workspace/vocabulary/vibe-coding-starter-api-go/CLAUDE.md
@/home/ubuntu/workspace/vocabulary/vibe-coding-starter-api-go/docs/code-generator.md

查看代码生成器使用方法，请生成一个 部门管理的前后端代码，在生成的代码基础上，再完成以下优化：
- 前端改为树形结构的部门管理
- 后端增加相应的api的接口
- 完成前端的多国语言翻译

====

启动前后端服务：

cd /home/ubuntu/workspace/vocabulary/vibe-coding-starter-api-go && go run cmd/server/main.go

cd /home/ubuntu/workspace/vocabulary/vibe-coding-starter-ui-antd && pnpm dev

进一步使用Playwright MCP进行测试：
- 使用分辨率：1920x1080
- 测试帐号/密码：admin/vibecoding

请验证`部门管理` CRUD功能是否正确，修复前后端的错误