# 门面模式（Facade Pattern）

    又叫作外观模式，提供了一个统一的接口，用来访问子系统中的一群接口。其主要特征是定义了一个高层接口，让子系统更容易使用
    
# 场景
    某在线商城, 推出了积分兑换礼品的功能，兑换礼品有几个步骤, 涉及到若干子系统:
    1。积分系统, 检查用户积分是否足够
    2。库存系统, 检查礼品是否有库存
    3。物流系统, 安排礼品发货并生成发货订单

# 设计
     GiftInfo: 礼品信息实体. 礼品也是一种库存物品.
     GiftExchangeRequest: 积分兑换礼品申请
     IGiftExchangeService: 积分兑换礼品服务, 该服务是一个Facade, 内部调用了多个子系统的服务
     IPointsService: 用户积分管理服务的接口
     IInventoryService: 库存管理服务的接口
     IShippingService: 物流下单服务的接口
     tMockGiftExchangeService: 积分兑换礼品服务的实现类
     tMockPointsService: 用户积分管理服务的实现类
     tMockInventoryService: 库存管理服务的实现类
     tMockShippingService: 物流下单服务的实现类
     
门面模式的优点

    （1）简化了调用过程，不用深入了解子系统，以防给子系统带来风险。
    （2）减少系统依赖，松散耦合。
    （3）更好地划分访问层次，提高了安全性。
    （4）遵循迪米特法则
门面模式的缺点

    （1）当增加子系统和扩展子系统行为时，可能容易带来未知风险。
    （2）不符合开闭原则。
    （3）某些情况下，可能违背单一职责原则