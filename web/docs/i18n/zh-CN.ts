export default {
  ProductionFeature: {
    title: '业界一流的文献管理',
    description: '高效地将"论文"转化为"知识"',
    features: {
      cloudSync: {
        title: '文献库免费云端同步',
        description: '跨设备同步您的文献库，确保随时随地都能访问您的研究资料。',
        imgUrl: '/images/library.png',
        iconUrl: '/images/icons/cloud-sync.svg',
      },
      folders: {
        title: '超级文件夹/标签，多级嵌套，组织清晰',
        description: '使用无限层级的文件夹和标签整理文献，构建您专属的知识树。',
        imgUrl: '/images/multi-level-folders.png',
        iconUrl: '/images/icons/folders.svg',
      },
      sorting: {
        title: '多维度排序和筛选',
        description: '基于各种元数据的强大排序和筛选功能，快速找到您需要的论文。',
        imgUrl: '/images/sorting-filtering.png',
        iconUrl: '/images/icons/sorting.svg',
      },
      citations: {
        title: '自定义引文、章节和卷号',
        description: '轻松管理和自定义引用详情，满足您的格式要求。',
        imgUrl: '/images/custom-citations.png',
        iconUrl: '/images/icons/citations.svg',
      },
    },
    buttons: {
      startReading: '开始阅读',
    },
  },
  AiNativeReaderFeature: {
    title: 'AI 原生阅读器',
    description: '利用 AI 简化阅读中繁琐的细节',
    features: {
      sideBySide: {
        title: '图表、文本，并列阅读',
        description: '毫不费力地比较图、表和文本，无需来回滚动，保持阅读连贯性。',
        imgUrl: '/images/ai-parser.png',
        iconUrl: '/images/icons/parser.svg',
      },
      jumpToReferences: {
        title: '一键跳转至参考文献',
        description: '通过即时跳转到您需要的参考文献，无缝地浏览您的文档。',
        imgUrl: '/images/ai-reference.png',
        iconUrl: '/images/icons/reference.svg',
      },
      noteTakingFormats: {
        title: '多种笔记格式：高亮、形状和文本框',
        description: '使用专为学术阅读定制的各种注释工具捕捉您的想法和见解。',
        imgUrl: '/images/ai-comment.png',
        iconUrl: '/images/icons/comment.svg',
      },
    },
    buttons: {
      startReading: '开始阅读',
    },
  },
  AIAssistedAnalysisFeature: {
    title: 'AI 辅助论文分析',
    description: '深入问题本质',
    features: {
      summaries: {
        title: '一键 AI 论文摘要',
        description: '通过简洁的AI生成摘要快速掌握论文核心思想，突出关键发现。',
        imgUrl: '/images/summaries.png',
        iconUrl: '/images/icons/summaries.svg',
      },
      savedQuestions: {
        title: '一键保存常问问题',
        description: '为每篇论文创建个性化问答库，加深理解并构建知识体系。',
        imgUrl: '/images/saved-asked-questions.png',
        iconUrl: '/images/icons/saved-asked-questions.svg',
      },
      instantAnswers: {
        title: '即时回答论文相关问题',
        description: '询问论文相关的任何问题，获得基于文档内容的即时、准确回答。',
        imgUrl: '/images/any-question.png',
        iconUrl: '/images/icons/any-question.svg',
      },
      followUpQuestions: {
        title: '无限追问',
        description: '通过持续的对话式问答体验深入探讨话题，每个回答都建立在前一个基础上。',
        imgUrl: '/images/unlimited-follow-up-questions.png',
        iconUrl: '/images/icons/unlimited-follow-up-questions.svg',
      },
      compareModels: {
        title: '对比最先进模型的结果',
        description: '利用多个领先AI模型的力量获得最佳见解，交叉验证答案准确性。',
        imgUrl: '/images/multi-models.png',
        iconUrl: '/images/icons/multi-models.svg',
      },
    },
    buttons: {
      startReading: '开始阅读',
    },
  },
  ProfessionalTranslationFeature: {
    title: '专业翻译与自定义术语库',
    description: '用强大的翻译工具克服研究中的语言障碍。',
    features: {
      fullText: {
        title: '全文翻译，沉浸式双语阅读体验',
        description: '用母语阅读论文，同时提供原文对照视图，确保完整理解。',
        imgUrl: '/images/full-translate.png',
        iconUrl: '/images/icons/full-translate.svg',
      },
      highlight: {
        title: '划词翻译：选中不懂的文本即可翻译',
        description: '无需离开页面或中断阅读流程，即可获得单词或短语的即时翻译。',
        imgUrl: '/images/translate.png',
        iconUrl: '/images/icons/translate.svg',
      },
      terminology: {
        title: '专注学术，可自定义术语库，科研词汇更准确',
        description: '使用专业的学术术语库确保翻译准确性，您可以根据自己的领域进行定制。',
        imgUrl: '/images/the-term.png',
        iconUrl: '/images/icons/the-term.svg',
      },
    },
    buttons: {
      startReading: '开始阅读',
    },
  },
  Pricing: {
    title: '套餐与价格',
    betaTitle: 'BETA期间，大幅优惠，BETA后恢复原价',
    subtitle: '选择适合您的完美计划',
    promo: 'chatgpt5/claude sonnet 4/claude sonnet3.7/gemini-pro-2.5，一个账户一网打尽',
    compareTitle: '套餐对比',
    month: '月',
    cascade: 'Cascade',
    free: {
      title: '免费版',
      price: '0',
      period: '每月',
      button: '开始',
      features: {
        trial: '专业版本试用',
        copilot: '使用 Ai Copilot',
        credits: '每月{0}个提示积分',
        betaCredits: 'beta期间每月包含{0}个提示积分',
        addonCredits: '无积分加油包',
        docmentCapacity: '文档云存储',
        fullTextTranslation: '全文翻译',
        aiTranslation: 'AI翻译',
      }
    },
    pro: {
      title: '专业版',
      price: '{0}{1}',
      betaPrice: '现价:',
      period: '每月',
      button: '选择此计划',
      features: {
        everything: '包含免费版的所有功能，以及：',
        credits: '每月{0}个提示积分',
        betaCredits: 'beta期间每月包含{0}个提示积分',
        copilot: '使用 Ai Copilot',
        addonCredits: '额外积分{0}{1}/{2}积分',
      }
    },
    enterprise: {
      title: '企业版',
      subtitle: '(超过200用户)',
      price: '联系我们',
      letsTalk: '联系我们'
    },
    table: {
      betaSign: 'BETA期间',
      promptCredits: '提示积分',
      creditsPerMonth: '{0}积分/月',
      addonCredits: '额外提示积分',
      addonCreditsPrice: '{0}{1}/{2}积分',
      features: '功能特性',
      aiCopilot: 'AI辅读',
      docs: '文档云存储',
      docMaxCapacity: '最大容量',
      docMaxStorageCapacity: '{0}MB',
      docMaxPageCount: '最大页数',
      docUploadMaxPageCount: '{0}页',
      docMaxSize: '最大文件大小',
      docUploadMaxSize: '每个文件 {0}MB',
      notes: "云笔记",
      fullTextTranslate: '全文翻译',
      fullTextTranslateCredits: '消耗积分',
      fullTextTranslateCreditsValue: '{0}积分',
      fullTextTranslateMaxPageCount: '最大页数',
      fullTextTranslateMaxPageCountValue: '{0}页',
      wordTranslate: '单词翻译',
      wordTranslateCredits: '消耗积分',
      wordTranslateCreditsValue: '{0}积分',
      aiTranslation: 'AI翻译',
      aiTranslationCredits: '消耗积分',
      aiTranslationCreditsValue: '{0}积分',
      ocrTranslate: 'OCR翻译',
      ocrTranslateCredits: '消耗积分',
      ocrTranslateCreditsValue: '{0}积分',
      startReading: '开始阅读',
      selectPlan: '选择方案',
      contactUs: '联系我们',
      costCredit: '费用积分'
    },
    faq: {
      title: '价格常见问题',
      stillHaveQuestions: '还有疑问？请联系我们的',
      support: '客服',
      team: '团队',
      questions: {
        promptCredit: {
          question: '什么是提示积分，它们如何被消耗？',
          answer: '每次向odoc发送消息时都会消耗一定的提示积分。使用的积分数量取决于模型和消息的长度。大多数交互使用1-2个积分。'
        },
        runOutCredits: {
          question: '当我用完积分后会发生什么？',
          answer: '当您用完积分后，您可以继续使用odoc的免费版功能。专业版用户可以随时购买额外积分，或等待每月积分刷新。'
        },
        studentPricing: {
          question: '你们提供学生价格吗？',
          answer: '是的！拥有有效.edu电子邮件地址的学生可以获得专业版计划50%的折扣。请携带您的学生证联系我们的客服团队获取折扣码。'
        },
        earlyAdopter: {
          question: '我的早期采用者价格会受影响吗？',
          answer: '不会，如果您在我们的早期采用者期间注册，您的特殊价格将被锁定，只要您保持订阅，价格就不会改变。'
        },
        autoRefills: {
          question: '附加积分充值是如何工作的？',
          answer: '专业用户可以在账户中购买附加积分包。'
        },
        subscription: {
          question: '为什么要订阅会员？',
          answer: '会员可优先体验AI功能，基于科研领域的学术大模型，可快速/准确理解复杂的学术文献和用户输入、并生成高质量回答及摘要，享受丰富图文信息的交互和沉浸式的问答体验，使用学术大模型能够帮助你打破传统学术的搜索限制，提升科研生产力，真正做到用AI轻松读论文。'
        },
        subscriptionEnd: {
          question: '会员到期后，会员权益会怎样？',
          answer: '会员到期后，会员权益会自动失效，您将无法继续使用会员功能。'
        },
        permanentPurchase: {
          question: '可以永久买断吗？',
          answer: 'odoc会员暂不支持永久买断制。若团队/企业有需要，请联系商务。'
        }
      }
    }
  },
  nav: {
    home: '首页',
    guide: '指南',
    membership: '会员',
    workbench: '我的文库',
  },
  layout: {
    userProfile: '个人资料',
    settings: '设置',
    logout: '退出登录',
    login: '登录',
  },
  ui: {
    login: '登录',
    logout: '退出登录',
    profile: '个人资料',
    settings: '设置',
  },
  footer: {
    product: {
      title: '产品',
      features: '功能特性',
      usage: '使用技巧',
      userAgreement: '用户协议',
      privacy: '隐私政策'
    },
    team: {
      title: '团队',
      aboutUs: '关于我们',
      updates: '更新动态',
      contactUs: '联系我们',
      joinUs: '加入我们'
    },
    social: {
      title: '社交媒体',
      one: 'Facebook',
      two: 'reddit',
      three: 'whatsapp'
    },
    copyright: '版权所有'
  },
  TargetUserScenarios: {
    title: '目标用户场景',
    tabs: {
      universityStudent: '大学生',
      universityProfessor: '大学教授',
      scientificResearcher: '科研人员',
      corporateRDEngineer: '企业研发工程师',
      lifelongLearner: '终身学习者',
    },
    content: {
      universityStudent: {
        title: '将杂乱无章、耗时数周的文献综述过程，转变为结构化、高效、深度思考的工作流。',
        description: [
          '<strong>快速筛选：</strong> Alex 将 100 篇论文批量导入软件。利用“一键AI摘要”功能，他快速阅读了每篇论文的AI生成摘要。仅用半天时间，他就将核心阅读清单缩小到 30 篇，大大节省了初步筛选的时间。',
          '<strong>深度阅读：</strong> 在阅读一篇核心论文时，他遇到了一个充满复杂公式的段落。他使用“针对特定选择提问”功能向 AI 提问：“能否用更简单的术语解释这个公式的物理意义？”AI 提供了清晰的解释，帮助他克服了理解障碍。',
          '<strong>知识连接：</strong> 每读完一篇论文，他都会创建“知识卡片”，总结核心思想、关键图表和自己的见解。最后，他将这些卡片排列在一个可视化白板上，绘制出不同算法的演进关系和优劣对比，使他的论文结构跃然纸上。',
          '<strong>标准化写作：</strong> 在撰写论文时，他使用 Word 插件一键插入引文，自动生成符合其大学要求的“标准格式参考文献”，满足了学校的要求。',
        ],
        scenario: {
          title: '场景',
          text: '他需要阅读近百篇关于最新AI算法的英文论文以完成文献综述，感到不堪重负和效率低下。',
        },
        role: {
          title: '角色',
          name: 'Alex',
          description: '一名计算机科学硕士生，正在准备毕业论文',
        },
        features: [
          '深度阅读',
          '标准化写作',
          '快速筛选',
          '知识连接',
        ],
        imageUrl: '/images/university-student.png',
      },
      universityProfessor: {
        title: '显著提高课程准备效率，为学生提供创新的、引导式的、互动式阅读体验，从单向讲座转向协作式知识建构。',
        description: [
          '<strong>课程准备：</strong> 李教授创建了一个名为"数字社会学课程"的共享文库。她添加了课程所需的数十篇核心论文和书籍章节。',
          '<strong>创建阅读指南：</strong> 对于每周的阅读材料，她首先自己阅读材料，高亮关键段落并添加笔记作为指导。然后，她让AI为每篇论文生成三个批判性思维问题，并将这些问题附加到她的笔记中。',
          '<strong>课堂互动：</strong> 她与学生分享文库。学生们可以提前阅读带注释的材料，并基于AI生成的问题进行小组讨论，显著提高了课堂的质量和互动性。',
        ],
        scenario: {
          title: '场景',
          text: '她需要准备课程材料，为学生创造引人入胜的学习体验，同时管理多个班级。',
        },
        role: {
          title: '角色',
          name: '李教授',
          description: '一名教授数字社会学的大学教授',
        },
        features: [
          '课程准备',
          '创建阅读指南',
          '课堂互动',
          '协作学习',
        ],
        imageUrl: '/images/university-professor.png',
      },
      scientificResearcher: {
        title: '将AI从简单的"问答工具"提升为激发科学灵感的"智能伙伴"，帮助研究人员在信息海洋中高效发现创新。',
        description: [
          '<strong>发现创新：</strong> 陈博士使用<strong>"对图表截图提问"</strong>功能。她选择数据图表并询问AI："除了作者的结论之外，这个图表还暗示了哪些其他潜在的相关性或异常？"',
          '<strong>持续探究：</strong> AI的回答激发了她的想法。她继续使用<strong>"无限追问"</strong>功能："基于你提到的异常，它可能与哪些已知的信号通路相关？"通过与AI的持续对话，她快速识别出一个之前被忽视的研究方向。',
          '<strong>项目提案：</strong> 她将AI辅助头脑风暴会议的见解整理成知识卡片，并将它们与其他几篇支持性论文链接起来，创建了一个逻辑严密、有充分支撑的结构化文库，为她的新研究项目提案奠定了坚实基础。',
        ],
        scenario: {
          title: '场景',
          text: '她需要发现新的研究方向，在分析复杂科学数据的同时制定创新的项目提案。',
        },
        role: {
          title: '角色',
          name: '陈博士',
          description: '一名国家级研究所的研究员',
        },
        features: [
          '发现创新',
          '持续探究',
          '项目提案',
          'AI辅助分析',
        ],
        imageUrl: '/images/scientific-researcher.png',
      },
      corporateRDEngineer: {
        title: '用集中化、可追溯的共享团队知识库，取代了旧有的混乱邮件往来和令人困惑的文档版本管理，显著加速团队的技术决策过程。',
        description: [
          '<strong>团队知识库：</strong> 张工程师创建了一个共享团队文库，并在其中添加了该技术的官方白皮书、核心研究论文和几篇关键技术博客文章。',
          '<strong>异步协作：</strong> 团队成员在阅读文档时，直接在PDF上添加他们的问题和见解作为笔记。这让每个人都能看到彼此的想法，促进了高效的异步沟通。',
          '<strong>技术对比：</strong> 在团队会议期间，张工程师打开软件并筛选所有成员的笔记来指导讨论。他还现场询问AI："请创建一个表格，从性能、成本和社区支持方面比较这项技术与我们当前的解决方案X。"AI快速生成了一个清晰简洁的表格。',
        ],
        scenario: {
          title: '场景',
          text: '他需要协调团队技术评估并做出明智决策，同时管理多个利益相关者和复杂的技术信息。',
        },
        role: {
          title: '角色',
          name: '张工程师',
          description: '一家科技公司的研发工程师',
        },
        features: [
          '团队知识库',
          '异步协作',
          '技术对比',
          '决策制定',
        ],
        imageUrl: '/images/corporate-r&d.png',
      },
      lifelongLearner: {
        title: '软件充当耐心、博学且全天候的"个人导师"，让深奥的科学知识变得易于理解，显著降低终身学习者进入高技术领域的门槛。',
        description: [
          '<strong>降低入门门槛：</strong> 她找到了史蒂芬·霍金的一篇经典论文，但发现其中充满了专业术语。她频繁使用划词翻译和针对性提问功能。例如，她高亮"事件视界"并询问："你能用外行人能理解的类比来解释这是什么吗？"',
          '<strong>构建知识体系：</strong> 每当她掌握一个新概念时，她就创建一张知识卡片，用自己的话总结。随着时间的推移，这些卡片形成了她个人的"黑洞知识图谱"。',
          '<strong>探索式学习：</strong> 在向AI询问一个概念后，她会继续问："基于这个，我接下来应该学习什么？"AI会推荐相关的经典论文或概念，引导她进行探索式学习之旅。',
        ],
        scenario: {
          title: '场景',
          text: '作为一名软件开发人员，她想在业余时间理解史蒂芬·霍金论文中的复杂天体物理学概念。',
        },
        role: {
          title: '角色',
          name: 'Sarah',
          description: '一位热衷于理解复杂科学概念的终身学习者',
        },
        features: [
          '降低入门门槛',
          '构建知识体系',
          '探索式学习',
          '个人辅导',
        ],
        imageUrl: '/images/lifelong-learner.png',
      },
    }
  },
  DifferWithChatGPT: {
    title: '我们与ChatGPT的区别',
    sections: {
      readingRetention: {
        title: '深度阅读：让知识真正留存',
        content: 'AI产品的频繁用户都知道，虽然AI提供了大量数据和比传统搜索引擎更快的答案，但存在一个问题。我们发现，由于这些信息获取得太容易，往往不会在我们的记忆中留存。真正的学习需要反复的输入和输出才能在人脑中编码。odoc的创建就是为了解决这个问题，帮助您"放慢"阅读节奏。通过高亮、注释和引导式阅读的沉浸式过程，我们的功能帮助将知识嵌入您的长期记忆中。'
      },
      professionalTool: {
        title: '专业的学术文献工具',
        content: '我们的功能集旨在覆盖整个研究生命周期：从文献收集 → 深度阅读 → 知识管理 → 论文写作 → 分享成果。我们自动获取和解析元数据（如作者、期刊、年份、DOI），确保准确性并为未来的引用和参考文献管理奠定坚实基础。界面和用户体验专为学术阅读的独特需求而定制，具有双窗格视图、参考链接和多文档比较等功能，使其区别于通用文件管理器或PDF阅读器。'
      },
      structuredLibraries: {
        title: '结构化文库，不仅仅是文件夹',
        content: '超越文件夹：这不仅仅是本地文件夹系统；它是一个多维数据库。用户可以使用文件夹、标签、评级、已读/未读状态、自定义字段等来组织、过滤和排序文献。深度交互：用户可以直接在原文上高亮、下划线和注释。这些笔记与文库深度集成，使其可搜索和可导出。系统化笔记：所有文档的所有笔记都可以在中央"笔记中心"中管理。用户可以按标签或按论文查看所有重要注释，而无需重新打开每个PDF。'
      },
      immersiveExperience: {
        title: '沉浸式"所见即所读"体验',
        content: '无干扰环境：我们提供干净、专注的阅读界面，去除所有不必要的杂乱。它支持多种主题（如护眼模式、暗黑模式）并适应不同的文档格式，让您纯粹专注于内容。流畅交互：所有工具栏和菜单都设计得直观——需要时出现，不需要时消失。翻译、记笔记或提问等操作都是无缝的，不会中断您的阅读流程。'
      },
      aiFeatures: {
        title: '一键AI摘要',
        content: '对于一篇长而陌生的论文，用户可以首先生成简洁摘要，快速掌握其背景、方法、发现和结论。'
      },
      unlimitedQuestions: {
        title: '交互式知识流：开启无限对话学习',
        content: '这开启了一种“苏格拉底式”的追问机制。 在AI给出解答后，您可以毫不费力地提出更深入、多层次的后续问题，形成一个自然、持续的知识对话流，让阅读过程转变为一场批判性思维的探索之旅。'
      },
      specificSelection: {
        title: '精准高亮解读：即时扫清微观障碍',
        content: '这是我们提供的极致精确的交互方式。 当遇到模棱两可的段落或复杂的措辞时，只需选中您想要理解的具体文本。您可以直接针对高亮部分提问，立即获得高度聚焦的、专为该片段定制的即时解答。'
      }
    }
  },
  UserTestimonials: {
    title: '用户评价',
    testimonials: {
      user1: {
        content: 'odoc是我最喜欢的文献阅读和管理工具。\n\nAI辅助阅读功能"非常出色"。导师布置的阅读任务，以前需要两三天完成，现在只需要半天。我可以直接向AI提问任何问题。这是一个巨大的生产力提升！\n我在每次小组会议中总是比实验室同学领先很多！',
        author: '@Researcher_TS8007A2',
        source: '来源：用户社区',
        avatar: 'user_1.svg'
      },
      user2: {
        content: '我试过很多类似的产品，但经过所有的搜索，我总是回到odoc。它在各个方面都提供了最流畅的体验。它在阅读和整理论文方面给了我很大帮助。我喜欢它的文献管理和AI写作辅助功能，这让我整个研究和写作过程更快更高效！',
        author: '@Researcher_TER4C0J2',
        source: '来源：用户社区',
        avatar: 'user_2.svg'
      },
      user3: {
        content: '有太多出色的功能无法一一列举！\n\n用文件夹和标签进行组织，让查找论文变得非常快速。\n\n划词翻译/全文翻译和AI阅读器让阅读论文变得轻松愉快。\n\n笔记和任务管理功能让我超级高效。\n感谢odoc！',
        author: '@Researcher_TSU1R562',
        source: '来源：用户社区',
        avatar: 'user_3.svg'
      },
      user4: {
        content: '我最喜欢的功能是翻译。作为一名国际学生，odoc的翻译功能显著缓解了用非母语进行研究的劣势。它帮助我在QS前15名大学保持一等荣誉学位，一篇论文接一篇论文。',
        author: '@momo',
        source: '来源：用户社区',
        avatar: 'user_4.svg'
      }
    }
  },
  OurPartners: {
    title: '合作伙伴',
    partners: {
      openai: {
        name: 'OpenAI',
        logo: 'openai.svg',
      },
      google: {
        name: 'Google',
        logo: 'google.svg',
      },
      anthropic: {
        name: 'Anthropic',
        logo: 'anthropic.svg',
      },
      microsoft: {
        name: 'Microsoft',
        logo: 'microsoft.svg',
      },
      harvard: {
        name: 'Harvard',
        logo: 'harvard.svg',
      },
      oxford: {
        name: 'Oxford',
        logo: 'oxford.svg',
      },
      yale: {
        name: 'Yale',
        logo: 'yale.svg',
      },
    },
  },
};
