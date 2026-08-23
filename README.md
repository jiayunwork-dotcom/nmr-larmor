# nmr-larmor

NMR 拉莫尔频率与化学位移计算工具。给定核素、外加磁场强度 B0 和化学位移（ppm），
计算该核在该场下的拉莫尔频率（Hz / MHz）、角频率，以及对应的绝对频率偏移（Hz）。
附带弛豫（T1/T2）、谱线线型、脉冲序列与磁场修正等扩展物理模型。

## 功能

- 拉莫尔共振：根据核素旋磁比与 B0 计算共振频率（ω = γ·B0）。
- 化学位移：将 ppm 相对位移换算为绝对频率偏移（δ_f = δ_ppm · f0 · 1e-6）。
- 弛豫模型：饱和恢复 / 反转恢复纵向弛豫、横向衰减、Ernst 角、CPMG 回波链。
- 谱线建模：洛伦兹 / 高斯线型、线宽与 T2 关系、ppm↔Hz 换算、耦合裂分。
- 脉冲序列：硬脉冲翻转角、Bloch 旋转、自旋回波、采样参数（采样间隔 / 采集时间）。
- 磁场修正：体磁化率去磁修正、参考标准位移、场扫描网格。

## 用法

### 命令行直接算

```bash
go run . -nucleus 1H -B0 7 -delta 1.5
```

### 跑内置算例

```bash
go run . -example example/h1-7t.json
```

### 启动网页控制台

```bash
go run . -http :8080
```

启动后在浏览器打开页面，可通过 `/api/larmor` 提交计算请求，或通过 `/api/scan`
做磁场扫描。

## 输入 / 输出

输入（JSON，对应 `LarmorRequest`）：

- `nucleus`：同位素符号，如 `1H`、`13C`、`31P`
- `B0`：主磁场强度，单位特斯拉（T），必须 > 0
- `delta_ppm`：相对参考（如 TMS）的化学位移，单位 ppm，可为 0
- `delta_unit`：化学位移单位，仅支持 `ppm`（留空等同 ppm）

输出（JSON，对应 `LarmorResult`）：

- `nucleus` / `name`：核素符号与名称
- `B0`：磁场强度
- `gamma_hz_per_t`：旋磁比（Hz/T）
- `f0_hz` / `f0_mhz`：拉莫尔频率（Hz / MHz）
- `omega_rad_s`：角频率（rad/s）
- `delta_f_hz`：绝对频率偏移（Hz）

## 构建与测试

```bash
go build ./...
go test ./...
```

核心算法位于 `internal/{nuclei,larmor,units,relaxation,spectrum,pulse,field}`，
HTTP 与 CLI 入口位于 `internal/httpapi` 与 `main.go`，仅做参数解析与连接。
