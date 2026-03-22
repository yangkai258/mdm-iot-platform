# MEMORY.md - agentcp 知识库

## 项目经验

### MDM 控制中台
- 项目路径: C:\Users\YKing\.openclaw\workspace\mdm-project
- 技术栈: Go+Gin后端 / Vue3+ArcoDesign前端
- 已产出: 会员管理需求文档、API文档

## 产品设计原则

1. **用户价值优先** - 每个功能问：用户真的需要吗？
2. **简单即美** - 能不复杂就不复杂
3. **数据驱动** - 决策基于数据，不是直觉
4. **迭代思维** - MVP先跑，快速验证

## API设计规范

### RESTful约定
- GET /resources - 列表
- GET /resources/:id - 详情
- POST /resources - 创建
- PUT /resources/:id - 更新
- DELETE /resources/:id - 删除

### 响应格式
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 常用工具

- API文档: Apifox / Swagger
- 原型设计: Figma
- 流程图: draw.io
- 需求管理: Notion / 飞书

## 学习方向

- [ ] 深入研究用户行为分析
- [ ] 学习数据分析方法
- [ ] 跟进竞品动态
- [ ] 提升交互设计能力

## 成长目标

- 短期: 完善会员管理模块文档
- 中期: 建立数据分析体系
- 长期: 成为领域专家

---

_保持好奇，保持批判，保持用户视角。_
