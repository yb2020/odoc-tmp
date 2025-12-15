# Membership Service - Order Module

## 积分支付订单状态机 (Credit Payment Status)

### 状态流转图 (Mermaid)

```mermaid
stateDiagram-v2
    direction LR
    [*] --> PENDING: 1. 创建订单

    PENDING --> AWAITING_CONFIRMATION: 2. 预扣积分成功
    PENDING --> FAILED: 2a. 预扣积分失败 (如积分不足)

    AWAITING_CONFIRMATION --> SUCCESS: 3. 业务确认 (Confirm)
    AWAITING_CONFIRMATION --> CANCELLED: 3a. 业务回滚 (Retrieve)

    subgraph "终态"
        SUCCESS
        FAILED
        CANCELLED
    end
```

### 流程说明

1.  **`[*] -> PENDING` (创建订单)**
    *   当用户发起一个积分支付请求时，系统创建一条新的支付记录，其初始状态为 `PENDING` (待处理)。

2.  **`PENDING -> AWAITING_CONFIRMATION` (预扣成功)**
    *   系统进入“准备”阶段，尝试从用户账户中**预先扣除**（或冻结）所需积分。
    *   如果操作成功，订单状态变为 `AWAITING_CONFIRMATION` (等待确认)，并等待上游业务的最终指令。

3.  **`PENDING -> FAILED` (预扣失败)**
    *   如果在预扣积分时失败（最常见的原因是用户积分不足），订单直接流转到 `FAILED` (支付失败)状态，流程结束。

4.  **`AWAITING_CONFIRMATION -> SUCCESS` (业务确认)**
    *   上游业务处理成功，调用 `ConfirmCreditFun` 接口。
    *   系统收到确认指令后，将订单状态更新为 `SUCCESS` (支付成功)。预扣的积分被正式划转，交易完成。

5.  **`AWAITING_CONFIRMATION -> CANCELLED` (业务回滚)**
    *   上游业务处理失败或超时，调用 `RetrieveCreditFun` 接口。
    *   系统收到回滚指令后，将预扣的积分**返还**给用户账户，并将订单状态更新为 `CANCELLED` (已取消)。

这个流程清晰地定义了二阶段提交（2PC）中的“准备” (`AWAITING_CONFIRMATION`)、“提交” (`SUCCESS`) 和“回滚” (`CANCELLED`) 三个关键环节，能够很好地支持你的业务需求。

## 会员订阅订单 (Subscription Order)

会员订阅订单用于用户购买或续费会员资格。这类订单通常涉及真实货币支付，其处理流程与支付网关（由 `pay` 服务封装）紧密集成。

### 订单状态流程

```mermaid
stateDiagram-v2
    direction LR
    [*] --> PENDING_PAYMENT: 1. 创建订单
    PENDING_PAYMENT --> COMPLETED: 2. 支付成功
    PENDING_PAYMENT --> FAILED: 2a. 支付失败
    PENDING_PAYMENT --> CANCELLED: 2b. 用户取消或超时

    subgraph "终态"
        COMPLETED
        FAILED
        CANCELLED
    end
```

### 流程说明

1.  **`[*] -> PENDING_PAYMENT` (待支付)**
    *   用户选择会员方案后，系统创建 `Order` 记录，状态为 `PENDING_PAYMENT`。
    *   系统调用 `pay` 服务获取支付凭证，并返回给前端。

2.  **`PENDING_PAYMENT -> COMPLETED` (已完成)**
    *   `pay` 服务通过 Webhook 通知支付成功。
    *   本服务验证通知，更新订单状态为 `COMPLETED`。
    *   **核心业务**：为用户开通或续期会员资格（更新 `UserMembership` 记录），并发放相应权益（如积分）。

3.  **`PENDING_PAYMENT -> FAILED` (支付失败)**
    *   支付网关返回明确的支付失败信息。订单状态更新为 `FAILED`。

4.  **`PENDING_PAYMENT -> CANCELLED` (已取消)**
    *   用户主动取消或支付超时。订单状态更新为 `CANCELLED`。

## 会员资格管理 (User Membership)

本模块负责管理用户的会员状态、等级和有效期。

### 核心设计

*   **自动续期**：免费会员资格采用自动续订模式。系统会为即将过期或已过期的免费会员自动续期，确保用户体验的连续性。
*   **新用户初始化**：新用户注册时，会自动获得一份初始的免费会员资格，无需等待后台任务触发。

## 积分退款 (Credit Refund)

为了保证账目清晰和数据不可变性，积分支付的退款流程遵循以下原则：

*   **独立退款记录**：退款操作不会修改原始的 `CreditPaymentRecord`。相反，系统会创建一条独立的 `CreditRefundRecord` 来记录退款详情。
*   **状态不变**：状态为 `SUCCESS` 的支付记录被视为历史凭证，其状态不会被更改。
*   **关联查询**：通过查询 `CreditRefundRecord` 表（可关联原始支付ID）来确定一笔支付是否已退款。

这种设计将支付和退款逻辑解耦，简化了对账和审计的复杂度。
