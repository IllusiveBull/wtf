# System Temperature Module

显示系统各硬件的温度信息，包括 CPU、GPU、SSD、内存和电池等。

## 系统要求

支持 macOS 系统（Intel 和 Apple Silicon）。

该模块使用 [iSMC](https://github.com/dkorunic/iSMC) Go 库直接读取 Apple SMC（系统管理控制器）传感器数据，无需额外安装命令行工具。

## 配置示例

```yaml
systemtemp:
  enabled: true
  position:
    top: 0
    left: 0
    height: 3
    width: 1
  refreshInterval: 5s
  showCPU: true
  showGPU: true
  showSSD: true
  showMem: true
  showBattery: true
```

## 配置参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `enabled` | boolean | `false` | 是否启用该模块 |
| `refreshInterval` | string | `"5s"` | 刷新间隔 |
| `showCPU` | boolean | `true` | 显示 CPU 温度（多核心平均值） |
| `showGPU` | boolean | `true` | 显示 GPU 温度 |
| `showSSD` | boolean | `true` | 显示 SSD/NAND 温度 |
| `showMem` | boolean | `true` | 显示内存温度 |
| `showBattery` | boolean | `true` | 显示电池温度 |

## 颜色说明

温度显示会根据温度值使用不同的颜色：

- **绿色**：< 50°C（正常）
- **黄色**：50-70°C（温暖）
- **橙色**：70-85°C（偏热）
- **红色**：> 85°C（过热）

## 支持的硬件

### Apple Silicon (M1/M2/M3 等)
- ✅ CPU Performance Core 温度
- ✅ CPU Efficiency Core 温度
- ✅ GPU 温度
- ✅ SSD/NAND 温度
- ✅ 内存温度
- ✅ 电池温度

### Intel Mac
- ✅ CPU 温度
- ✅ GPU 温度（如果有独立显卡）
- ✅ SSD 温度
- ✅ 内存温度
- ✅ 电池温度（笔记本）

## 技术细节

该模块使用 [iSMC](https://github.com/dkorunic/iSMC) Go 库，该库能够：
- 直接访问 Apple SMC 传感器
- 支持 Intel 和 Apple Silicon Mac
- 自动检测和解码温度传感器
- 支持 HID 传感器（Apple Silicon）

## 注意事项

1. **权限**：某些系统可能需要管理员权限才能访问 SMC
2. **传感器可用性**：不同 Mac 型号可用的传感器可能不同
3. **刷新间隔**：建议不要设置得太短（推荐 5 秒或更长），以避免频繁读取传感器
4. **温度精度**：温度读数精度取决于硬件传感器

## 故障排除

### 无法获取温度信息

可能的原因：

1. **权限不足**：尝试使用 `sudo` 运行程序
2. **SMC 访问被限制**：检查系统安全设置
3. **不支持的硬件**：某些 Hackintosh 或虚拟机可能不支持

### CPU 温度显示异常

- CPU 温度为多个核心的平均值
- Apple Silicon Mac 有性能核心和效率核心，模块会读取所有核心温度并计算平均值

### 某些温度不显示

- 不是所有 Mac 都有所有类型的传感器
- 可以在配置中关闭不需要的温度显示
- 台式机 Mac（如 Mac mini、iMac）可能没有电池温度

## 示例输出

```
System Temperature
------------------
CPU:       54.2°C
GPU:       48.3°C
SSD:       41.5°C
Memory:    47.8°C
Battery:   32.1°C
```

## 相关项目

- [iSMC](https://github.com/dkorunic/iSMC) - Apple SMC CLI 工具和 Go 库
- [stats](https://github.com/exelban/stats) - macOS 系统监控应用
- [iStats](https://github.com/Chris911/iStats) - Ruby Gem Mac 系统状态工具
