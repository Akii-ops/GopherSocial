package db

import (
	"backend/internal/store"
	"context"
	"fmt"
	"log"
	"math/rand"
)

var blogUsernames = []string{
	"Emily", "Liam", "Sophia", "Noah", "Olivia",
	"Ava", "Elijah", "Mia", "James", "Charlotte",
	"Benjamin", "Amelia", "Lucas", "Harper", "Mason",
	"Evelyn", "Ethan", "Abigail", "Logan", "Ella",
	"Alexander", "Grace", "Jacob", "Chloe", "Michael",
	"Penelope", "Daniel", "Lily", "Henry", "Aria",
	"Jackson", "Zoey", "Sebastian", "Riley", "Aiden",
	"Nora", "Matthew", "Hazel", "Samuel", "Ellie",
	"David", "Violet", "Joseph", "Luna", "Carter",
	"Stella", "Owen", "Lucy", "Wyatt", "Claire",
}

var blogTitles = []string{
	"AI绘画革命:Midjourney实战指南",
	"从零搭建React18性能优化体系",
	"元宇宙社交:Horizon Worlds初体验",
	"深度解读《三体》中的Dark Forest理论",
	"Go语言并发编程:Channel高级技巧",
	"Vintage风格装修:老宅Loft改造实录",
	"Web3.0入门:ENS域名投资指南",
	"咖啡地图:手冲咖啡的变量控制Variables",
	"Flutter跨平台开发:iOS/Android实战",
	"碳中和时代:ESG投资逻辑解析",
	"暗黑系摄影:Low-key Lighting布光教学",
	"智能家居:HomeKit自动化场景设计",
	"剧本杀创作:Plot Twist设计方法论",
	"植物养护指南:多肉Succulent度夏技巧",
	"古典乐赏析:贝多芬Moonlight Sonata",
	"极简主义:Capsule Wardrobe构建法则",
	"深度学习实战:PyTorch模型优化技巧",
	"城市探索:Abandoned Factory摄影纪实",
	"健康管理:Biohacking基础实践",
	"独立开发:从MVP到Product Hunt发布",
}

var blogContents = []string{
	"详解MidJourney基础命令与常用参数组合，通过prompt工程实现特定艺术风格输出", // 对应AI绘画标题
	"通过代码分割/懒加载优化首屏性能，分析React18并发渲染的核心机制",
	"对比Horizon Worlds与VRChat的社交体系差异，实测虚拟身份创建流程",
	"解读黑暗森林法则的博弈论基础，推演费米悖论的可能解决方案",
	"深入讲解select语句与context包在goroutine生命周期管理中的实践",
	"记录老厂房改造中如何保留工业元素，平衡现代极简设计语言",
	"解析ENS域名估值模型，演示在OpenSea平台进行域名交易的完整流程",
	"量化研磨度/水温/注水速度对风味的影响，建立可复用的冲煮参数矩阵",
	"实现跨平台状态共享方案，解决iOS/Android系统特性差异导致的UI适配问题",
	"拆解环境社会治理指标，分析新能源产业在ESG评级中的权重变化",
	"演示单灯布光技巧，利用伦勃朗光塑造戏剧化的人物肖像效果",
	"构建基于地理围栏的自动化场景，实现离家自动关断高耗能设备",
	"设计多线叙事结构，通过红鲱鱼技巧制造推理反转效果",
	"针对不同品种制定配土方案，建立夏季控水与通风的标准化流程",
	"解析第三乐章的情感表达，比较不同钢琴家演奏的速度处理差异",
	"建立33件基础单品数据库，制定季节性胶囊衣橱的搭配算法",
	"应用混合精度训练与知识蒸馏技术，提升模型推理效率30%以上",
	"记录工业遗址摄影中的构图思维，危险环境拍摄的安全防护要点",
	"实践冷暴露与间歇性断食方案，量化监测心率变异性改善幅度",
	"制定产品上线Checklist，涵盖ASO优化与媒体包准备的完整流程",
}

var blogTags = [][]string{
	{"AI绘画", "数字艺术", "工具教程", "AIGC"},
	{"前端开发", "React", "性能优化", "Web应用"},
	{"元宇宙", "虚拟现实", "社交网络", "科技趋势"},
	{"科幻文学", "博弈论", "物理学", "书评"},
	{"Go语言", "并发编程", "系统设计", "后端开发"},
	{"室内设计", "旧房改造", "工业风", "生活美学"},
	{"区块链", "Web3.0", "数字资产", "投资指南"},
	{"咖啡文化", "手冲技巧", "精品咖啡", "生活仪式感"},
	{"移动开发", "跨平台", "Flutter", "UI适配"},
	{"可持续发展", "ESG投资", "碳中和", "商业分析"},
	{"摄影艺术", "人像摄影", "布光技巧", "暗黑系"},
	{"智能家居", "HomeKit", "物联网", "自动化"},
	{"剧本创作", "推理游戏", "叙事设计", "社交娱乐"},
	{"园艺养护", "多肉植物", "度夏技巧", "绿植"},
	{"古典音乐", "钢琴曲", "音乐赏析", "贝多芬"},
	{"极简主义", "胶囊衣橱", "生活方式", "断舍离"},
	{"深度学习", "PyTorch", "模型优化", "AI工程"},
	{"城市探索", "废墟摄影", "工业遗产", "纪实"},
	{"健康科技", "生物黑客", "自我量化", "健康管理"},
	{"独立开发", "产品设计", "创业", "MVP"},
}

var blogComments = []string{
	"Prompt工程部分讲得太到位了！终于搞明白seed参数怎么用了",
	"React并发模式的实际性能提升有测试数据支持吗？",
	"虚拟社交的身份切换功能如果能支持跨平台就好了",
	"黑暗森林理论与囚徒困境的结合角度很新颖！",
	"Go的context超时控制在实际微服务中确实关键",
	"工业风吊顶的管线处理方案可以再详细些吗？",
	"ENS域名投资的风险评估部分建议补充案例",
	"手冲变量矩阵图太实用了，已打印贴在咖啡台！",
	"跨平台状态管理方案能否用于电商类APP？",
	"ESG指标在传统制造业的应用有参考资料推荐吗？",
	"伦勃朗光教学视频的机位角度能展示下吗？",
	"地理围栏触发条件在多层公寓的精度如何？",
	"红鲱鱼技巧在剧本杀中的占比多少比较合适？",
	"配土方案能具体到颗粒土品牌吗？",
	"不同演奏版本的频谱分析图很有说服力！",
	"胶囊衣橱的算法逻辑是否有开源项目参考？",
	"知识蒸馏的教师模型选择有什么讲究？",
	"工业遗址拍摄的安全装备清单很专业！",
	"心率变异性检测用什么设备比较准？",
	"Product Hunt上线前的邮件预热模板能分享吗？",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}

	}

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("seeding complete")

	return

}

func generateUsers(num int) []*store.User {

	users := make([]*store.User, num)

	for i := 0; i < num; i++ {

		users[i] = &store.User{
			Username: blogUsernames[i%len(blogUsernames)] + fmt.Sprintf("%d", i),
			Email:    blogUsernames[i%len(blogUsernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}

	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Version: 0,
			Title:   blogTitles[rand.Intn(len(blogTitles))],
			Content: blogContents[rand.Intn(len(blogContents))],
			Tags:    blogTags[rand.Intn(len(blogTags))],
		}
	}

	return posts

}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		post := posts[rand.Intn(len(posts))]

		comments[i] = &store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: blogComments[rand.Intn(len(blogComments))],
		}
	}

	return comments
}
