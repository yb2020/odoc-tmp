# 支付子系统 (Pay Module)

## 一、核心目标

*   **可靠支付处理**: 为应用提供稳定、安全的支付功能。
*   **Stripe 优先支持**: 初期完美集成 Stripe 作为主要的支付网关。
*   **高扩展性**: 架构设计应易于未来接入更多支付渠道（如 PayPal、支付宝、微信支付等）。
*   **松耦合集成**: 与项目中其他模块（如订单模块、会员模块）保持低耦合，方便独立开发和维护。

## 二、当前实现状态

本支付子系统已实现以下核心功能：

*   **支付渠道抽象**: 通过 `PaymentProvider` 接口和 `PaymentProviderFactory` 支持多种支付渠道的集成，目前已实现 `StripeProvider`。
*   **支付核心服务**: `PaymentService` 编排支付创建、Webhook 处理、退款、状态查询等核心逻辑。
*   **API 接口**: 提供 RESTful API 用于发起支付、处理 Stripe Webhook、查询支付状态、发起退款及查询用户支付记录。
*   **数据持久化**: `PaymentRecordDAO` 结合 `GormBaseDAO` 管理 `PaymentRecord` 数据的存储与检索。
*   **标准化模型**: `PaymentRecord` 模型遵循项目规范，嵌入 `BaseModel`，并包含支付相关的详细字段和状态常量。

## 三、主要组件及职责

1.  **API 层 (`services/pay/api/payment_api.go`)**
    *   **职责**: 作为支付模块的 HTTP 入口，处理客户端请求和支付渠道回调。
    *   **主要端点**:
        *   `POST /payments/prepay`: 发起支付请求。
        *   `POST /payments/stripe/webhook`: 处理 Stripe Webhook 事件。
        *   `GET /payments/{payment_id}/status`: 查询特定支付的状态。
        *   `POST /payments/{payment_id}/refund`: 为特定支付发起退款。
        *   `GET /payments/user`: 查询指定用户的支付记录列表。
    *   **技术栈**: 使用 `gin-gonic/gin` 构建，响应格式遵循 `pkg/response` 标准。

2.  **服务层 (`services/pay/service/payment_service.go`)**
    *   **`PaymentService`**: 包含支付业务的核心逻辑和编排。
    *   **主要职责**:
        *   `CreatePayment`: 验证参数，创建内部 `PaymentRecord`，并调用相应支付渠道的 `CreateCharge` 方法。
        *   `HandleWebhook`: 验证并处理来自支付渠道的 Webhook 事件，更新支付记录状态。
        *   `CreateRefund`: 处理退款请求，调用支付渠道的退款接口，并更新支付记录。
        *   `GetPaymentStatus`: 查询并返回支付记录的当前状态。
        *   `GetPaymentsByUserId`: 根据用户ID和可选状态筛选支付记录。
    *   **依赖**: `PaymentRecordDAO` (数据持久化), `PaymentProviderFactory` (获取支付渠道实例)。

3.  **支付渠道抽象层 (`services/pay/provider/`)**
    *   **`PaymentProvider` 接口 (`payment_provider.go`)**: 定义支付渠道的统一操作标准。
        *   `GetName() string`: 返回支付渠道名称。
        *   `CreateCharge(params CreateChargeParams) (*CreateChargeResult, error)`: 创建支付/扣款。
        *   `HandleWebhook(payload []byte, signature string) (*WebhookEvent, error)`: 解析并验证 Webhook 数据。
        *   `CreateRefund(params CreateRefundParams) (*RefundResult, error)`: 发起退款。
        *   `GetChargeStatus(providerTxId string) (string, error)`: 查询支付渠道的交易状态。
    *   **`StripeProvider` (`stripe_provider.go`)**: `PaymentProvider` 接口的 Stripe 实现。
        *   封装 Stripe Go SDK (v82) 的调用逻辑。
        *   处理 Stripe API 密钥和 Webhook 签名密钥的配置。
    *   **`PaymentProviderFactory` (`provider_factory.go`)**: 注册和获取 `PaymentProvider` 实例的工厂。
        *   `RegisterProvider(provider PaymentProvider)`: 注册支付渠道提供者。
        *   `GetProvider(name string) (PaymentProvider, error)`: 根据名称获取提供者实例。

4.  **数据访问层 (DAO) (`services/pay/dao/payment_record_dao.go`)**
    *   **`PaymentRecordDAO`**: 负责 `PaymentRecord` 数据的持久化和读取。
    *   继承自 `pkg/dao.GormBaseDAO[model.PaymentRecord]`，复用通用的 CRUD 方法 (如 `Save`, `FindById`)。
    *   提供特定查询方法，例如：
        *   `GetByOrderID(ctx context.Context, orderID string) (*model.PaymentRecord, error)`
        *   `GetByProviderTxId(ctx context.Context, providerTxId string) (*model.PaymentRecord, error)`
        *   `UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error`
        *   `GetByUserIdAndStatus(ctx context.Context, userId string, status string, page, size int) ([]*model.PaymentRecord, int64, error)`

5.  **数据模型 (`services/pay/model/payment_record.go`)**
    *   **`PaymentRecord`**: 支付主记录，用于持久化支付交易的详细信息。
        *   **基础模型**: 嵌入 `pkg/model.BaseModel`，自动继承 `Id` (int64, 雪花ID), `IsDeleted`, `CreatedAt`, `UpdatedAt`, `CreatorId`, `ModifierId` 等字段和 GORM 钩子。
        *   **核心业务字段**:
            *   `UserId` (string): 发起支付的用户ID。
            *   `OrderId` (string): 关联的业务订单ID。
            *   `Amount` (int64): 支付金额，以最小货币单位表示（例如：分）。
            *   `Currency` (string): ISO 4217 货币代码 (例如 "CNY", "USD")。
            *   `Status` (string): 支付状态 (如 `PaymentStatusPending`, `PaymentStatusSucceeded`)。
            *   `Channel` (string): 支付渠道 (如 `PaymentChannelStripe`)。
            *   `ProviderTxId` (string): 支付渠道返回的唯一交易ID。
        *   **详细信息与元数据**:
            *   `Description` (string): 支付描述。
            *   `Metadata` (map[string]interface{}): 自定义元数据，数据库类型为 `jsonb`。
            *   `PaymentMethodDetails` (map[string]interface{}): 支付方式详情，数据库类型为 `jsonb`。
            *   `ProviderErrorCode` (string): 支付渠道返回的错误码。
            *   `ProviderErrorMessage` (string): 支付渠道返回的错误信息。
            *   `PaidAt` (*time.Time): 支付成功的时间。
        *   **数据库表名**: 通过 `TableName() string` 方法指定为 `t_payment_record`。
        *   **常量**: 定义了 `PaymentStatus...` 和 `PaymentChannel...` 系列常量。

6.  **配置管理 (Conceptual)**
    *   **职责**: 集中管理支付渠道相关的配置信息 (如 API 密钥、Webhook Secret)。
    *   Stripe 相关配置: `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET` (后端用), `STRIPE_PUBLISHABLE_KEY` (前端用)。
    *   **实践**: 这些配置应通过环境变量或安全的配置文件管理，不应硬编码在代码中。

## 四、核心交互流程 (以 Stripe 支付为例 - 当前实现)

1.  **客户端 -> API 层**: 客户端（如前端应用）调用 `POST /payments/prepay`，提供用户信息、订单信息、金额、币种、支付渠道 (e.g., "STRIPE") 及支付方式ID (Stripe `PaymentMethodID`)。
2.  **API 层 -> 服务层**: `PaymentAPI` 将请求参数传递给 `PaymentService.CreatePayment`。
3.  **服务层 (`PaymentService.CreatePayment`)**:
    a.  创建 `PaymentRecord` 实例，初始状态为 `PaymentStatusPending`。
    b.  调用 `PaymentRecordDAO.Save` 持久化支付记录。
    c.  通过 `PaymentProviderFactory` 获取 `StripeProvider`。
    d.  调用 `StripeProvider.CreateCharge`，传递支付参数。
4.  **支付渠道层 (`StripeProvider.CreateCharge`)**:
    a.  使用 Stripe Go SDK 调用 Stripe API 创建 `PaymentIntent` (通常包含 `confirm: true` 以尝试立即支付)。
    b.  返回 `CreateChargeResult`，包含 `ProviderTxId` (Stripe `PaymentIntentID`), `ClientSecret` (如果需要前端额外操作), 支付状态等。
5.  **服务层 -> API 层**: `PaymentService` 返回 `CreatePaymentResult` 给 `PaymentAPI`。
6.  **API 层 -> 客户端**: `PaymentAPI` 将结果（如 `ClientSecret`）返回给客户端，客户端可使用 Stripe.js 处理后续步骤（如 3D Secure 验证）。
7.  **Stripe Webhook -> API 层**: 当支付状态在 Stripe 端发生变化 (如 `payment_intent.succeeded`, `payment_intent.failed`)，Stripe 会向配置的 `POST /payments/stripe/webhook` 端点发送事件。
8.  **API 层 (`PaymentAPI.HandleStripeWebhook`)**:
    a.  读取请求体和 `Stripe-Signature` 头。
    b.  调用 `PaymentService.HandleWebhook`，传递渠道名 ("STRIPE")、请求体和签名。
9.  **服务层 (`PaymentService.HandleWebhook`)**:
    a.  获取 `StripeProvider`。
    b.  调用 `StripeProvider.HandleWebhook` 验证签名并解析事件。
    c.  根据事件类型更新 `PaymentRecord` 状态 (通过 `PaymentRecordDAO.UpdateFields`)。
    d.  (未来可扩展) 触发后续业务逻辑（如通知订单模块）。

## 五、关键设计考虑点 (与当前实现相关)

*   **幂等性**: Webhook 处理逻辑应设计为幂等的。当前实现通过查询支付记录状态并在更新前进行检查，部分实现了幂等性。更完善的方案可考虑记录已处理的 Webhook 事件 ID。
*   **安全性**:
    *   Stripe API 密钥和 Webhook Secret 需安全配置。
    *   Stripe Webhook 签名验证已在 `StripeProvider.HandleWebhook` 中实现。
*   **异步处理**: Webhook 的处理是异步的，确保状态更新的最终一致性。
*   **错误处理与日志**: 各层均包含错误处理和日志记录 (使用 `pkg/logging`)。
*   **可扩展性**: `PaymentProvider` 接口和 `PaymentProviderFactory` 为接入新支付渠道提供了基础。

此 README 描述了 `pay` 子系统的当前实现状态和核心架构，为团队成员理解模块功能和未来开发提供了指引。

## 六、前端与后端交互流程 (Vue 与 Stripe)

此流程描述了典型的基于Stripe Elements的前端（以Vue为例）与后端支付服务的交互：

**前端 (Vue) - 步骤 1: 初始化与信息收集**
1.  **加载Stripe.js**: 在Vue应用中安全地加载Stripe.js。
2.  **初始化Stripe Elements**: 使用Publishable Key初始化Stripe实例，并创建和挂载Elements（如Card Element）到DOM中，供用户输入支付信息。
3.  **收集订单信息**: 用户确认支付意愿后，前端收集必要的订单详情（如金额、商品信息等，这些信息可能部分来自后端）。

**前端 (Vue) - 步骤 2: 创建支付方式并发起支付**
4.  **创建PaymentMethod**: 当用户提交支付表单时，Vue组件调用 `stripe.createPaymentMethod({type: 'card', card: cardElement})` (或其他支付方式对应的 Elements) 从Stripe Elements中安全地收集支付信息并创建 `PaymentMethodID`。
5.  **调用后端发起支付接口**: Vue组件向后端API `POST /payments/prepay` 发送请求，请求体中包含上一步获取的 `PaymentMethodID` 以及其他必要的订单详情（如用户ID、订单ID、金额、币种等）。

**后端 - 步骤 3: 处理支付请求 (如本文档 "四、核心交互流程" 所述)**
6.  **创建PaymentRecord**: 后端服务接收到请求后，首先创建一条内部的 `PaymentRecord`，初始状态通常为 `Pending`。
7.  **创建Stripe PaymentIntent**: 调用 `StripeProvider.CreateCharge` 方法。此方法内部会使用Stripe SDK向Stripe API发起请求创建 `PaymentIntent`。
    *   通常会设置 `confirm: true` 尝试立即确认支付。
    *   传递从前端获取的 `PaymentMethodID`。
    *   如果支付方式需要用户在银行页面进行额外验证（如3D Secure），则必须提供 `ReturnURL`。

**前端 (Vue) - 步骤 4: 处理后端响应与后续操作**
8.  **接收后端响应**: 前端接收来自 `/payments/prepay` 的响应。此响应通常包含：
    *   `paymentIntentId`: Stripe PaymentIntent的ID。
    *   `clientSecret`: PaymentIntent的客户端密钥。
    *   `status`: PaymentIntent的当前状态 (e.g., `succeeded`, `requires_action`, `requires_payment_method`)。
9.  **处理PaymentIntent状态**:
    *   **支付成功 (`succeeded`)**: 如果状态是 `succeeded`，则支付已完成。前端更新UI，向用户显示成功信息。
    *   **需要额外操作 (`requires_action`)**: 如果状态是 `requires_action` (例如，3D Secure验证)，前端需要使用返回的 `clientSecret` 调用 `stripe.confirmCardPayment(clientSecret)` (或其他特定支付方式的确认函数，如 `stripe.confirmIdealPayment` 等)。这将引导用户完成必要的验证步骤（可能会重定向到银行页面）。
    *   **支付失败或需要新支付方式 (`requires_payment_method`, `failed`)**: 向用户显示错误信息，并可能允许用户尝试使用新的支付方式。
10. **处理 `confirmCardPayment` 结果**:
    *   调用 `stripe.confirmCardPayment` 后，会返回一个包含更新后 `PaymentIntent` 对象或错误的对象。
    *   前端根据这个结果更新UI，显示支付成功、失败或仍在处理中的状态。

**前端 (Vue) - 步骤 5: 处理重定向 (如果发生)**
11. **用户重定向**: 如果在步骤 9 或 10 中用户被重定向到第三方页面（如银行的3D安全验证页面）并完成操作，Stripe会将用户重定向回应用中预设的 `ReturnURL`。
12. **在ReturnURL页面检索状态**: `ReturnURL` 对应的Vue组件或页面加载时，应从URL查询参数中获取 `payment_intent_client_secret` (Stripe会自动附加)。
13. 使用此 `client_secret` 调用 `stripe.retrievePaymentIntent(client_secret)` 从Stripe获取最新的 `PaymentIntent` 状态。
14. 根据检索到的最终状态更新UI，向用户展示支付结果。

**后端 (Webhook) - 步骤 6: 最终状态确认 (如本文档 "四、核心交互流程" 第7-9点所述)**
15. **接收Webhook事件**: Stripe会异步发送Webhook事件 (如 `payment_intent.succeeded`, `payment_intent.payment_failed`, `charge.succeeded` 等) 到后端配置的 `/payments/stripe/webhook` 端点。
16. **处理Webhook**: 后端验证Webhook签名，解析事件，并据此更新内部 `PaymentRecord` 的最终状态。**这是确保支付状态与Stripe侧最终一致的最可靠方式，不应仅依赖前端的反馈。**

**总结**:
这种设计将敏感的支付信息处理（如信用卡号）隔离在由Stripe提供的安全UI组件 (Stripe Elements) 和Stripe的服务器中。
*   **前端职责**: 主要负责用户交互、通过Stripe.js安全地收集支付信息（生成`PaymentMethodID`）、处理需要前端介入的支付步骤（如3D Secure）、以及与自身后端API的通信。
*   **后端职责**: 负责核心的支付业务逻辑处理、与Stripe服务器进行安全通信（创建和管理`PaymentIntent`、处理Webhook）、数据持久化以及维护支付状态的权威记录。

## 七、支付模块与订单模块交互流程

支付模块与业务订单模块的顺畅交互对于保障整体业务流程至关重要。以下是典型的交互模式：

**核心思想：** 订单模块管理用户订单的生命周期（用户希望购买或使用的服务/商品），支付模块则负责处理该订单的金融交易。两者需有效通信以保持状态同步。

**典型交互流程：**

1.  **订单创建与待支付 (订单模块 -> 支付模块)**
    *   用户在应用中完成购物车结算或服务选择。
    *   **订单模块** 在其数据库中创建订单记录，包含唯一的 `OrderID` 及初始状态（如 `PendingPayment` - 待支付）。
    *   订单模块调用支付模块的支付发起接口（例如，本系统中的 `POST /payments/prepay`），传递 `OrderID`、`Amount` (金额)、`Currency` (币种)、`UserID` 等关键信息。

2.  **支付处理 (支付模块)**
    *   **支付模块** 接收到请求后，在其数据库中创建 `PaymentRecord`，并与传入的 `OrderID` 关联，初始状态通常为 `PaymentStatusPending` (处理中)。
    *   支付模块通过选定的支付渠道（如 Stripe）执行实际的支付操作，如本文档 “四、核心交互流程” 和 “六、前端与后端交互流程” 所述。

3.  **支付状态更新与通知 (支付模块 -> 订单模块)**
    *   **Webhook (主要且最可靠的方式):** 当支付渠道（如 Stripe）完成支付处理后（例如 `payment_intent.succeeded` 或 `payment_intent.payment_failed`），会向支付模块配置的 Webhook 端点 (例如 `/payments/stripe/webhook`) 发送事件通知。
    *   支付模块的 `PaymentService.HandleWebhook` 方法处理此事件，更新内部 `PaymentRecord` 的状态。
    *   **通知订单模块:** `PaymentRecord` 状态更新后，支付模块需要将支付结果（成功或失败）通知给订单模块。这是实现 “*(未来可扩展) 触发后续业务逻辑（如通知订单模块）*” (见 “四、核心交互流程” 第 9.d 点) 的关键步骤。
        *   **通知方式建议：**
            *   **直接 API 调用：** 支付服务调用订单模块暴露的内部 API（例如 `POST /internal/orders/{orderId}/payment-update`）。
            *   **消息队列/事件总线 (推荐，实现解耦):** 支付服务发布一个事件 (如 `PaymentSucceededEvent`、`PaymentFailedEvent`，包含 `OrderID` 和支付详情) 到消息队列 (如 RabbitMQ, Kafka)。订单模块订阅这些事件并相应处理。这符合 “松耦合集成” 的设计目标。

4.  **订单履行 (订单模块)**
    *   **订单模块** 收到支付成功的通知后，更新其对应订单的状态（例如 `Paid` - 已支付, `Processing` - 处理中, `Active` - 已激活等）。
    *   随后触发后续业务逻辑，如发送确认邮件、通知发货、开通服务权限等。
    *   若支付失败，订单模块也应更新订单状态（例如 `PaymentFailed` - 支付失败），并可能通知用户或提供重试选项。

**关键连接点：**
*   `PaymentRecord.OrderId` 字段是连接支付记录与具体业务订单的核心纽带。

**总结系统交互：**

*   **订单模块 (外部服务):**
    1.  创建订单，生成 `OrderID`。
    2.  调用 `pay` 服务的 `/payments/prepay` 接口，传递 `OrderID`、`Amount` 等。
*   **`pay` 支付服务:**
    1.  创建与 `OrderID` 关联的 `PaymentRecord`。
    2.  通过 Stripe (或其他渠道) 处理支付。
    3.  通过 Webhook 接收支付结果，更新 `PaymentRecord.Status`。
    4.  **主动通知订单模块** 支付结果 (携带 `OrderID`)。
*   **订单模块:**
    1.  接收来自 `pay` 服务的支付结果通知。
    2.  更新自身订单状态，并执行后续业务流程。

此交互模式确保了订单状态能够准确反映支付结果，是保障业务流程顺畅运行的基础。

## 八、Vue前端下单到支付的整体流程与页面设计

本章节概述了Vue前端从用户下单到完成支付的典型流程，以及建议设计的相关页面或组件。

**核心流程概述：**

用户选择商品 -> 创建订单 -> 选择支付方式 -> 输入支付信息 -> 处理支付（可能涉及银行验证） -> 显示支付结果 -> 后端异步确认最终状态。

**建议设计的Vue页面/组件及其职责：**

1.  **商品/服务列表页 (`ProductListPage.vue`)**
    *   **职责**: 展示可供购买的商品或服务。
    *   **交互**: 用户浏览并选择商品，可跳转到详情页或直接加入购物车。

2.  **商品详情页 (`ProductDetailPage.vue`) (可选)**
    *   **职责**: 展示单个商品/服务的详细信息。
    *   **交互**: 用户选择规格、数量，点击“加入购物车”或“立即购买”。

3.  **购物车页 (`ShoppingCartPage.vue`)**
    *   **职责**: 展示用户已添加到购物车的商品列表。
    *   **交互**: 用户修改商品数量、移除商品、查看总价，点击“去结算”。

4.  **订单确认/结算页 (`CheckoutPage.vue`)**
    *   **职责**: 显示订单摘要（商品、价格、收货地址等），用户确认订单信息。
    *   **后端交互**: 用户确认后，前端调用**后端订单模块API**创建订单，获取`OrderID`。
    *   **交互**: 确认后，引导进入支付环节。

5.  **支付方式选择与处理页/组件 (`PaymentPage.vue` 或嵌入 `CheckoutPage.vue`)**
    *   **职责**: 允许用户选择支付渠道（如Stripe、支付宝等），并展示相应支付输入界面。
    *   **后端交互 (以Stripe为例)**:
        1.  用户提交支付信息，前端通过Stripe.js获取`paymentMethodId`。
        2.  前端调用后端支付服务的 `POST /payments/prepay` 接口，传递`orderId`、`amount`、`currency`、`channel`及渠道特定参数（如Stripe的`paymentMethodId`和`returnUrl`）。
        3.  前端接收`prepay`接口响应，其中包含如`clientSecret` (Stripe)和`status`等信息，用于指导下一步操作。
    *   **前端处理 `prepay` 响应 (Stripe示例)**:
        *   若`status: 'succeeded'`：支付成功，跳转到支付成功页。
        *   若`status: 'requires_action'` (如3D Secure)：前端使用`clientSecret`调用`stripe.confirmCardPayment(clientSecret)`，引导用户完成验证。
        *   若支付失败：显示错误信息，允许重试。

6.  **支付重定向处理页 (`PaymentReturnPage.vue`)**
    *   **职责**: 作为Stripe等支付渠道在用户完成外部验证（如3D Secure）后的`returnUrl`。
    *   页面加载时，从URL查询参数中获取`payment_intent_client_secret` (Stripe会自动附加)。
    *   使用此`client_secret`调用`stripe.retrievePaymentIntent(client_secret)`获取最新的支付状态。
    *   **交互**: 根据最终状态导航到相应的支付结果页。

7.  **支付结果页 (`PaymentStatusPage.vue`)**
    *   **支付成功视图 (`PaymentSuccessView.vue`)**: 显示支付成功信息、订单号、金额等，提供“查看订单”或“返回首页”链接。
    *   **支付失败视图 (`PaymentFailureView.vue`)**: 显示支付失败信息、原因（若有）、订单号，提供“重新支付”或“联系客服”链接。

8.  **用户订单列表页 (`OrderListPage.vue`)**
    *   **职责**: 展示用户的历史订单及其状态。订单的最终权威状态依赖后端Webhook的异步确认。

9.  **订单详情页 (`OrderDetailPage.vue`)**
    *   **职责**: 展示特定订单的详细信息，包括支付状态和详情。

**整体前端流程串联 (以Stripe信用卡支付为例):**

1.  用户在商品页/购物车页选择商品，进入订单确认/结算页。
2.  在订单确认/结算页，用户确认信息，点击“提交订单并支付”。前端调用订单服务API创建订单，获得`OrderID`。
3.  进入支付方式选择与处理环节。用户选择信用卡，输入信息，点击“确认支付”。
    *   前端获取`paymentMethodId` (via Stripe.js)。
    *   前端调用后端 `POST /payments/prepay`。
4.  前端处理`prepay`接口响应：
    *   **直接成功**: 跳转支付成功页。
    *   **需额外操作 (e.g., 3D Secure)**:
        *   调用`stripe.confirmCardPayment(clientSecret)`。
        *   若无重定向，根据结果跳转支付结果页。
        *   若有重定向（如银行页面），用户完成后会返回到应用的`returnUrl` (进入步骤5)。
    *   **失败**: 提示错误。
5.  用户被重定向到支付重定向处理页 (对应`returnUrl`)：
    *   页面提取`payment_intent_client_secret`，调用`stripe.retrievePaymentIntent()`获取最终状态。
    *   跳转相应的支付结果页。
6.  在支付结果页向用户展示信息。

**重要提示:** 即使用户在前端看到支付成功，订单的最终权威状态依赖于后端通过支付渠道的Webhook异步接收并处理支付事件，然后更新订单模块中的订单状态。

## 九、后端 Stripe Checkout 模式集成 (方式二：独立服务)

在 Stripe Checkout 模式下，若采用保持 `StripeCheckoutService` 独立，并让 API Handler 直接调用它的方式，整合思路如下：

这种方式下，`StripeCheckoutService` 专注于处理 Stripe Checkout 的所有逻辑（创建 Session、处理 Checkout 相关的 Webhook），而您现有的 `PaymentService` 和 `StripeProvider` 继续负责它们原有的职责（例如基于 PaymentIntents 的支付流程）。

**文件结构设想:**

```
services/pay/
├── api/
│   ├── payment_api.go         // 现有API，可能包含 /payments/stripe-webhook
│   └── checkout_api.go        // (新增) 处理 /checkout/... 相关路由
├── service/
│   ├── payment_service.go     // 现有服务
│   └── stripe_checkout_service.go // (新增) Stripe Checkout 专属服务
├── provider/
│   ├── payment_provider.go    // 接口
│   ├── stripe_provider.go     // Stripe PaymentIntent 等实现
│   └── provider_factory.go
├── dao/
│   └── payment_record_dao.go
└── model/
    └── payment_record.go
```

**1. `StripeCheckoutConfig` (在 `stripe_checkout_service.go` 或共享的配置包中)**

```go
package service // 或 config

import "your_project_path/services/pay/dao" // 假设DAO在此

type StripeCheckoutConfig struct {
    SecretKey     string
    WebhookSecret string // 用于验证 Checkout Webhook 的签名密钥
    SuccessURL    string
    CancelURL     string
    Currency      string // 默认货币，例如 "usd", "cny"
}

// StripeCheckoutService 结构体
type StripeCheckoutService struct {
    cfg             StripeCheckoutConfig
    paymentRecordDAO dao.PaymentRecordDAO // 依赖注入 PaymentRecordDAO
    // 可能还有 logger 等其他依赖
}

// NewStripeCheckoutService 构造函数
func NewStripeCheckoutService(cfg StripeCheckoutConfig, prDAO dao.PaymentRecordDAO) *StripeCheckoutService {
    // stripe.Key = cfg.SecretKey // 初始化 Stripe Go SDK 的全局 API Key (在服务启动时更合适)
    return &StripeCheckoutService{
        cfg:             cfg,
        paymentRecordDAO: prDAO,
    }
}
```
*   **注意:** `stripe.Key` 的初始化最好在应用启动时进行一次，而不是在每次创建服务实例时。

**2. `stripe_checkout_service.go` 的核心方法**

```go
package service

import (
    "context"
    "encoding/json"
    "fmt"
    "io"      // 使用 io.ReadAll 替代 ioutil.ReadAll
    "net/http"
    "time"

    "github.com/stripe/stripe-go/v76" // 确保版本号与您项目一致
    "github.com/stripe/stripe-go/v76/checkout/session"
    "github.com/stripe/stripe-go/v76/webhook"

    "your_project_path/services/pay/model"
    "your_project_path/services/pay/dao"
)

// CreateCheckoutSessionParams 定义创建会话的输入参数
type CreateCheckoutSessionParams struct {
    ItemName  string `json:"itemName"`
    ItemPrice int64  `json:"itemPrice"` // 单位：分
    Quantity  int64  `json:"quantity"`
    OrderID   string `json:"orderId"`   // 内部订单ID
    UserID    string `json:"userId"`    // 用户ID (如果需要记录到PaymentRecord)
}

// CreateCheckoutSessionResult 定义创建会话的输出结果
type CreateCheckoutSessionResult struct {
    SessionID string `json:"sessionId"`
}

// CreateCheckoutSession 创建 Stripe Checkout Session
func (s *StripeCheckoutService) CreateCheckoutSession(ctx context.Context, params CreateCheckoutSessionParams) (*CreateCheckoutSessionResult, error) {
    stripe.Key = s.cfg.SecretKey // 确保Stripe Key已设置 (全局或每次调用前)
    stripeParams := &stripe.CheckoutSessionParams{
        Mode: stripe.String(string(stripe.CheckoutSessionModePayment)), // 或 CheckoutSessionModeSubscription
        LineItems: []*stripe.CheckoutSessionLineItemParams{
            {
                PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
                    Currency: stripe.String(s.cfg.Currency), // 从配置获取默认货币
                    ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
                        Name: stripe.String(params.ItemName),
                    },
                    UnitAmount: stripe.Int64(params.ItemPrice),
                },
                Quantity: stripe.Int64(params.Quantity),
            },
        },
        SuccessURL:        stripe.String(s.cfg.SuccessURL),
        CancelURL:         stripe.String(s.cfg.CancelURL),
        ClientReferenceID: stripe.String(params.OrderID), // 关键：关联内部订单
        Metadata: map[string]string{"user_id": params.UserID}, // 可选，用于传递额外信息
    }

    sess, err := session.New(stripeParams)
    if err != nil {
        return nil, fmt.Errorf("failed to create checkout session: %w", err)
    }

    return &CreateCheckoutSessionResult{SessionID: sess.ID}, nil
}

// HandleCheckoutWebhook 处理 Stripe Checkout Webhook 事件
func (s *StripeCheckoutService) HandleCheckoutWebhook(ctx context.Context, payload []byte, signatureHeader string) error {
    stripe.Key = s.cfg.SecretKey // 确保Stripe Key已设置
    event, err := webhook.ConstructEvent(payload, signatureHeader, s.cfg.WebhookSecret)
    if err != nil {
        return fmt.Errorf("webhook signature verification failed: %w", err)
    }

    switch event.Type {
    case "checkout.session.completed":
        var checkoutSession stripe.CheckoutSession
        err := json.Unmarshal(event.Data.Raw, &checkoutSession)
        if err != nil {
            return fmt.Errorf("error parsing webhook json for checkout.session.completed: %w", err)
        }
        return s.processCheckoutSessionCompleted(ctx, &checkoutSession)

    // case "checkout.session.async_payment_succeeded":
    // case "checkout.session.async_payment_failed":
    // 对于订阅模式，还需要处理以下事件：
    // case "customer.subscription.created":
    // case "customer.subscription.updated":
    // case "customer.subscription.deleted":
    // case "invoice.payment_succeeded": // 订阅续费成功
    // case "invoice.payment_failed":    // 订阅续费失败
    
    default:
        // log.Printf("Unhandled Stripe event type: %s", event.Type)
    }

    return nil
}

// processCheckoutSessionCompleted 处理 checkout.session.completed 事件的业务逻辑
func (s *StripeCheckoutService) processCheckoutSessionCompleted(ctx context.Context, cs *stripe.CheckoutSession) error {
    if cs.ClientReferenceID == "" {
        return fmt.Errorf("client_reference_id is empty in checkout.session.completed")
    }
    // 对于银行卡等同步支付，PaymentStatus 应该是 paid
    // 对于某些异步支付方式，可能需要等待 async_payment_succeeded 事件
    if cs.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
        // log.Printf("Checkout session %s not paid, status: %s. Awaiting async event or check configuration.", cs.ID, cs.PaymentStatus)
        // 根据业务决定是否处理或记录警告。如果只处理同步支付成功，可以返回错误或nil。
        return fmt.Errorf("checkout session %s payment status is %s, expected 'paid' for direct processing", cs.ID, cs.PaymentStatus)
    }

    // Idempotency: 检查是否已为该 OrderID 处理过成功的支付
    // existingRecord, err := s.paymentRecordDAO.GetByOrderID(ctx, cs.ClientReferenceID)
    // if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { return err }
    // if existingRecord != nil && existingRecord.Status == model.PaymentStatusSucceeded {
    //    log.Printf("Payment for order %s already processed.", cs.ClientReferenceID)
    //    return nil //幂等性：已处理
    // }

    paymentRecord := &model.PaymentRecord{
        UserId:       cs.Metadata["user_id"], // 假设 UserId 通过 Metadata 传递
        OrderId:      cs.ClientReferenceID,
        Amount:       cs.AmountTotal,
        Currency:     string(cs.Currency),
        Status:       model.PaymentStatusSucceeded, // 映射为内部成功状态
        Channel:      "STRIPE_CHECKOUT", // 定义一个常量 model.PaymentChannelStripeCheckout
        ProviderTxId: cs.ID,             // Stripe Checkout Session ID (cs_xxx)
        Description:  fmt.Sprintf("Stripe Checkout for Order %s", cs.ClientReferenceID),
        PaidAt:       stripe.Int64(time.Now().Unix()), // 支付完成时间, cs.Created 可能是会话创建时间
    }

    if err := s.paymentRecordDAO.Save(ctx, paymentRecord); err != nil {
        return fmt.Errorf("failed to save payment record: %w", err)
    }

    // (重要) 触发后续业务逻辑: 更新订单状态、通知用户等。
    // log.Printf("Successfully processed checkout.session.completed for OrderID: %s", cs.ClientReferenceID)
    return nil
}
```

**3. `checkout_api.go` (API Handler)**

```go
package api

import (
    "io"
    "net/http"

    "github.com/gin-gonic/gin"
    "your_project_path/services/pay/service"
)

type CheckoutAPI struct {
    checkoutService *service.StripeCheckoutService
}

func NewCheckoutAPI(cs *service.StripeCheckoutService) *CheckoutAPI {
    return &CheckoutAPI{checkoutService: cs}
}

func (a *CheckoutAPI) RegisterCheckoutRoutes(router *gin.RouterGroup) {
    checkoutGroup := router.Group("/checkout") // 例如 /api/v1/checkout
    checkoutGroup.POST("/create-session", a.createCheckoutSessionHandler)
    // 为 Checkout Webhook 单独设置端点和 Secret 更简单直接
    checkoutGroup.POST("/stripe-webhook", a.handleStripeCheckoutWebhookHandler) 
}

func (a *CheckoutAPI) createCheckoutSessionHandler(c *gin.Context) {
    var req service.CreateCheckoutSessionParams
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
        return
    }

    // userID := c.GetString("userID") // 从 context 或 token 中获取 UserID
    // req.UserID = userID // 并赋值给 req.UserID

    result, err := a.checkoutService.CreateCheckoutSession(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session: " + err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}

func (a *CheckoutAPI) handleStripeCheckoutWebhookHandler(c *gin.Context) {
    const MaxBodyBytes = int64(65536) // 64KB
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
    payload, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.Status(http.StatusServiceUnavailable)
        return
    }

    signatureHeader := c.GetHeader("Stripe-Signature")
    if err := a.checkoutService.HandleCheckoutWebhook(c.Request.Context(), payload, signatureHeader); err != nil {
        // log.Printf("Error handling Stripe Checkout webhook: %v", err)
        if _, ok := err.(stripe.SignatureVerificationError); ok { // 更具体的错误检查
            c.Status(http.StatusBadRequest) // 签名错误
        } else {
            c.Status(http.StatusInternalServerError) // 其他内部错误
        }
        return
    }
    c.Status(http.StatusOK)
}
```

**4. 服务和路由的初始化 (例如在 `cmd/server/main.go` 或类似地方)**

```go
// import (
//  payAPI "your_project_path/services/pay/api"
//  payDAO "your_project_path/services/pay/dao"
//  payService "your_project_path/services/pay/service"
//  "github.com/spf13/viper" // 示例配置库
//  "github.com/stripe/stripe-go/v76" // Stripe SDK
// )

func main_example_setup() { // 示例函数名，避免冲突
    // ... (加载配置 viper.GetXXX, 初始化数据库 db *gorm.DB, 初始化 Gin router *gin.Engine 等)

    // 1. 加载 Stripe Checkout 配置
    stripeCfg := payService.StripeCheckoutConfig{
        SecretKey:     viper.GetString("stripe.secretKey"),
        WebhookSecret: viper.GetString("stripe.checkoutWebhookSecret"), // Checkout 专属 Webhook Secret
        SuccessURL:    viper.GetString("stripe.checkoutSuccessUrl"),
        CancelURL:     viper.GetString("stripe.checkoutCancelUrl"),
        Currency:      viper.GetString("stripe.currency"),
    }

    // 全局初始化 Stripe API Key (应用启动时一次即可)
    stripe.Key = stripeCfg.SecretKey

    // 2. 初始化 DAO
    // paymentRecordDAO := payDAO.NewPaymentRecordDAO(db)

    // 3. 初始化 StripeCheckoutService
    // stripeCheckoutSvc := payService.NewStripeCheckoutService(stripeCfg, paymentRecordDAO)

    // 4. 初始化 CheckoutAPI
    // checkoutAPI := payAPI.NewCheckoutAPI(stripeCheckoutSvc)

    // 5. 注册路由 (假设 apiV1Group 是 /api/v1 或类似的基础路由组)
    // apiV1Group := router.Group("/api/v1") 
    // checkoutAPI.RegisterCheckoutRoutes(apiV1Group.Group("/payments")) // 例如 /api/v1/payments/checkout/create-session
    // 或者
    // checkoutAPI.RegisterCheckoutRoutes(apiV1Group) // 例如 /api/v1/checkout/create-session

    // ... (启动 Gin 服务器)
}
```

**关键整合点和决策:**

*   **Webhook 处理策略:** 推荐为 Checkout 事件使用独立的 Webhook 端点和 Secret (如示例中 `/checkout/stripe-webhook`)，这样更清晰简单。
*   **`PaymentRecord` 的创建/更新:** `StripeCheckoutService` 中的 `processCheckoutSessionCompleted` 负责。务必实现幂等性处理。
*   **错误处理和日志:** 在所有关键路径添加。
*   **配置管理:** 安全管理所有密钥和URL。
*   **依赖注入:** 服务和处理器通过构造函数注入依赖。
*   **Stripe API Key:** 全局设置 `stripe.Key = cfg.SecretKey` 一次，通常在应用启动时。如果服务需要处理多个Stripe账户，则需要在每次API调用前设置 `stripe.Key` 或使用 `stripe.APIKey` 参数。

这种方式下，`StripeCheckoutService` 成为一个专门处理Stripe Checkout流程的内聚模块，与您现有的基于PaymentIntents的支付逻辑可以并行存在，互不干扰，同时又能共享底层的 `PaymentRecordDAO` 和数据模型。

## 八、Stripe支付模式与Webhook事件处理对比表

根据Stripe官方文档，以下是一次性支付模式（payment mode）和订阅支付模式（subscription mode）需要处理的webhook事件对比：

### 支付模式对比表

| 事件类型 | 一次性支付模式 (payment) | 订阅支付模式 (subscription) | 说明 |
|---------|----------------------|--------------------------|------|
| **checkout.session.completed** | ✅ 必须处理 | ✅ 必须处理 | 初始支付成功，两种模式都需要处理 |
| **checkout.session.async_payment_succeeded** | ✅ 必须处理 | ✅ 必须处理 | 异步支付方式成功，两种模式都需要处理 |
| **checkout.session.async_payment_failed** | ✅ 必须处理 | ✅ 必须处理 | 异步支付方式失败，两种模式都需要处理 |
| **customer.subscription.created** | ❌ 不适用 | ✅ 必须处理 | 订阅创建成功，仅订阅模式需要处理 |
| **customer.subscription.updated** | ❌ 不适用 | ✅ 必须处理 | 订阅状态变更，仅订阅模式需要处理 |
| **customer.subscription.deleted** | ❌ 不适用 | ✅ 必须处理 | 订阅取消，仅订阅模式需要处理 |
| **invoice.payment_succeeded** | ❌ 不适用 | ✅ 必须处理 | 周期性付款成功，仅订阅模式需要处理 |
| **invoice.payment_failed** | ❌ 不适用 | ✅ 必须处理 | 周期性付款失败，仅订阅模式需要处理 |

### 支付流程对比

#### 一次性支付模式流程
1. 创建 Checkout Session (mode=payment)
2. 用户完成支付
3. 触发 **checkout.session.completed** 事件
4. 处理订单完成逻辑

#### 订阅支付模式流程
1. 创建 Checkout Session (mode=subscription)
2. 用户完成首次支付
3. 触发 **checkout.session.completed** 事件
4. 触发 **customer.subscription.created** 事件
5. 处理订阅创建逻辑
6. 周期性付款时：
   - 触发 **invoice.payment_succeeded** 事件（成功）
   - 或触发 **invoice.payment_failed** 事件（失败）
7. 订阅变更时触发 **customer.subscription.updated** 事件
8. 订阅取消时触发 **customer.subscription.deleted** 事件

### 关键差异

1. **初始支付**：两种模式都使用 checkout.session.completed 事件
2. **后续付款**：
   - 一次性支付：没有后续付款
   - 订阅支付：通过 invoice.payment_succeeded/failed 事件处理
3. **生命周期管理**：
   - 一次性支付：支付完成后结束
   - 订阅支付：需要管理完整的订阅生命周期（创建、更新、取消）
