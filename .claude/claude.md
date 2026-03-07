# Claude 项目指令

## 文件路径规范

在处理文件时，始终使用相对路径而非绝对路径。

### 示例

**错误：**
```
D:/desktop/sentinel/Sentinel/web/src/pages/subscriptions/index.tsx
```

**正确：**
```
web/src/pages/subscriptions/index.tsx
```

### 原因

- 相对路径更简洁易读
- 避免暴露用户的文件系统结构
- 提高代码的可移植性
- 减少路径冲突和错误

## 项目根目录

项目根目录为：`Sentinel/`

所有相对路径都基于此目录。