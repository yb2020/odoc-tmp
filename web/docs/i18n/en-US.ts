export default {
  ProductionFeature: {
    title: 'Best-in-class Literature Management',
    description: 'Effectively turn "papers" into "knowledge"',
    features: {
      cloudSync: {
        title: 'Free cloud synchronization for your literature library',
        description: 'Keep your library synced across all devices, ensuring you have access to your research anytime, anywhere.',
        imgUrl: '/images/library.png',
        iconUrl: '/images/icons/cloud-sync.svg',
      },
      folders: {
        title: 'Super folders/tags with multi-level nesting for crystal-clear organization',
        description: 'Organize your literature with infinite-level folders and tags to build your own knowledge tree.',
        imgUrl: '/images/multi-level-folders.png',
        iconUrl: '/images/icons/folders.svg',
      },
      sorting: {
        title: 'Multi-dimensional sorting and filtering',
        description: 'Quickly find the papers you need with powerful sorting and filtering options based on various metadata.',
        imgUrl: '/images/sorting-filtering.png',
        iconUrl: '/images/icons/sorting.svg',
      },
      citations: {
        title: 'Customizable citations, sections, and volume numbers',
        description: 'Easily manage and customize citation details to fit your formatting requirements.',
        imgUrl: '/images/custom-citations.png',
        iconUrl: '/images/icons/citations.svg',
      },
    },
    buttons: {
      startReading: 'Start Reading',
    },
  },
  AiNativeReaderFeature: {
    title: 'AI-Native Reader',
    description: 'AI like an expert next to you',
    features: {
      sideBySide: {
        title: 'View figures and tables side-by-side with the text',
        description: 'Effortlessly compare figures, tables, and text without losing context or scrolling back and forth.',
        imgUrl: '/images/ai-parser.png',
        iconUrl: '/images/icons/parser.svg',
      },
      jumpToReferences: {
        title: 'Jump to cited references with a single click',
        description: 'Navigate your documents seamlessly by instantly jumping to the references you need.',
        imgUrl: '/images/ai-reference.png',
        iconUrl: '/images/icons/reference.svg',
      },
      noteTakingFormats: {
        title: 'Multiple note-taking formats: highlights, shapes, and text boxes',
        description: 'Capture your thoughts and insights with a variety of annotation tools tailored for academic reading.',
        imgUrl: '/images/ai-comment.png',
        iconUrl: '/images/icons/comment.svg',
      },
    },
    buttons: {
      startReading: 'Start Reading',
    },
  },
  AIAssistedAnalysisFeature: {
    title: 'AI-Assisted Paper Analysis',
    description: 'AI-assisted Paper Understanding Comprehensive Explanation',
    features: {
      summaries: {
        title: 'One-click AI-powered paper summaries',
        description: 'Quickly grasp the core ideas of a paper with concise, AI-generated summaries that highlight key findings.',
        imgUrl: '/images/summaries.png',
        iconUrl: '/images/icons/summaries.svg',
      },
      savedQuestions: {
        title: 'Save frequently asked questions with one click',
        description: 'Create a personalized Q&A library for each paper to deepen your understanding and build knowledge.',
        imgUrl: '/images/saved-asked-questions.png',
        iconUrl: '/images/icons/saved-asked-questions.svg',
      },
      instantAnswers: {
        title: 'Instant answers to any question about the paper',
        description: 'Ask anything about the paper and get immediate, context-aware answers backed by the document content.',
        imgUrl: '/images/any-question.png',
        iconUrl: '/images/icons/any-question.svg',
      },
      followUpQuestions: {
        title: 'Unlimited follow-up questions',
        description: 'Dive deeper into topics with a continuous, conversational Q&A experience that builds on previous answers.',
        imgUrl: '/images/unlimited-follow-up-questions.png',
        iconUrl: '/images/icons/unlimited-follow-up-questions.svg',
      },
      compareModels: {
        title: 'Compare results with the most advanced models',
        description: 'Leverage the power of multiple leading AI models to get the best possible insights and cross-validate answers.',
        imgUrl: '/images/multi-models.png',
        iconUrl: '/images/icons/multi-models.svg',
      },
    },
    buttons: {
      startReading: 'Start Reading',
    },
  },
  ProfessionalTranslationFeature: {
    title: 'Professional translation with customizable dictionary',
    description: 'Overcome language barriers in your research with powerful translation tools.',
    features: {
      fullText: {
        title: 'Full-text translation for an immersive, bilingual reading experience',
        description: 'Read papers in your native language with a side-by-side view of the original text for complete understanding.',
        imgUrl: '/images/full-translate.png',
        iconUrl: '/images/icons/full-translate.svg',
      },
      highlight: {
        title: 'Highlight-to-translate: just select the text you don\'t understand',
        description: 'Get instant translations for words or phrases without leaving the page or interrupting your reading flow.',
        imgUrl: '/images/translate.png',
        iconUrl: '/images/icons/translate.svg',
      },
      terminology: {
        title: 'Academically focused with a customizable terminology base for more accurate scientific vocabulary',
        description: 'Ensure translation accuracy with a specialized academic terminology base that you can customize for your field.',
        imgUrl: '/images/the-term.png',
        iconUrl: '/images/icons/the-term.svg',
      },
    },
    buttons: {
      startReading: 'Start Reading',
    },
  },
  Pricing: {
    title: 'Plans and Pricing',
    betaTitle: 'BETA discounts available. Full price resumes after BETA',
    subtitle: 'Choose the right plan to support your research journey',
    promo: 'ChatGPT-5, Claude Sonnet 4, Claude Sonnet 3.7, and Gemini Pro 2.5 — all under a single account',
    compareTitle: 'Compare Plans',
    month: 'month',
    vibeReading: 'Vibe Reading',
    free: {
      title: 'Free',
      price: '{0}{1}',
      period: 'per month',
      button: 'Starting Now',
      features: {
        trial: 'One-time Pro Trial',
        copilot: 'Use Ai Copilot in Trial',
        credits: 'Includes {0} prompt credits per month',
        betaCredits: 'Includes {0} prompt credits per month during the beta period',
        addonCredits: 'No additional credits available',
        docmentCapacity: 'Document cloud storage',
        fullTextTranslation: 'Full-text translation',
        aiTranslation: 'Ai Translation',
      },
    },
    pro: {
      title: 'Pro',
      price: '{0}{1}',
      betaPrice: 'Price:',
      period: 'per month',
      button: 'Select Plan',
      features: {
        everything: 'Everything in Free, plus:',
        copilot: 'AI Copilot included in Pro',
        credits: '{0} prompt credits/month',
        betaCredits: '{0} prompt credits/month during the beta period',
        addonCredits: 'Purchase additional credits: {0}{1} for {2} credits',
      },
    },
    enterprise: {
      title: 'Enterprise',
      subtitle: '(more than 200 users)',
      price: "Let's talk",
      letsTalk: "Let's talk",
    },
    table: {
      betaSign: 'BETA discounts',
      promptCredits: 'Prompt credits',
      creditsPerMonth: '{0} credits/mo',
      addonCredits: 'Add-on prompt credits',
      addonCreditsPrice: '{0}{1} for {2} credits',
      features: 'Features',
      aiCopilot: 'AI Copilot',
      docs: 'Document cloud storage',
      docMaxCapacity: 'Max capacity',
      docMaxStorageCapacity: '{0}MB',
      docMaxPageCount: 'Max page count',
      docUploadMaxPageCount: '{0} pages',
      docMaxSize: 'Max size',
      docUploadMaxSize: 'Maximum file size: {0}MB',
      notes: "Cloud Notes",
      fullTextTranslate: 'Full Text Translate',
      fullTextTranslateCredits: 'Credit Cost',
      fullTextTranslateCreditsValue: '{0} credits',
      fullTextTranslateMaxPageCount: 'Max page count',
      fullTextTranslateMaxPageCountValue: '{0} pages',
      wordTranslate: 'Word Translate',
      wordTranslateCredits: 'Credit Cost',
      wordTranslateCreditsValue: '{0} credits',
      aiTranslation: 'AI Translation',
      aiTranslationCredits: 'Credit Cost',
      aiTranslationCreditsValue: '{0} credits',
      ocrTranslate: 'OCR Translate',
      ocrTranslateCredits: 'Credit Cost',
      ocrTranslateCreditsValue: '{0} credits',
      startReading: 'Start Reading',
      selectPlan: 'Select Plan',
      contactUs: 'Contact Us',
      costCredit: 'Credit Cost'
    },
    faq: {
      title: 'Pricing FAQ',
      stillHaveQuestions: 'Still have questions? Contact our',
      support: 'support',
      team: 'team',
      questions: {
        promptCredit: {
          question: 'What are prompt credits and how are they consumed?',
          answer: 'A prompt credit is consumed each time you send a message to Cascade. The number of credits used depends on the model and length of the message. Most interactions use 1-2 credits.'
        },
        runOutCredits: {
          question: 'What happens when I run out of credits?',
          answer: 'When you run out of credits, you can continue using Cascade\'s free features. Pro users can purchase additional credits at any time or wait for monthly credit refresh.'
        },
        studentPricing: {
          question: 'Do you offer student pricing?',
          answer: 'Yes! Students with a valid .edu email address can get 50% off the Pro plan. Contact our support team with your student ID to get a discount code.'
        },
        earlyAdopter: {
          question: 'Will my early-adopter discount remain the same?',
          answer: 'No, if you signed up during our early adopter period, your special pricing is locked in and won\'t change as long as you maintain your subscription.'
        },
        autoRefills: {
          question: 'How do automatic credit refills work?',
          answer: 'Pro users can enable automatic credit refills in their account settings. When your credits drop below 50, we\'ll automatically add 250 credits to your account for $10.'
        },
        subscription: {
          question: 'What are the benefits of subscribing?',
          answer: 'Members can experience AI features first, based on academic fields, they can quickly/accurately understand complex academic literature and user input, and generate high-quality answers and summaries, enjoy an interactive and immersive Q&A experience with rich graphic and text information, using academic models can help you break the traditional academic search limit, improve research productivity, and truly use AI to read papers easily.'
        },
        subscriptionEnd: {
          question: 'What will happen to my membership benefits after my membership expires?',
          answer: 'After your membership expires, your membership benefits will automatically expire, and you will not be able to continue using membership features.'
        },
        permanentPurchase: {
          question: 'Do you offer a lifetime license?',
          answer: 'odoc membership does not support permanent licensing. If your team or company needs it, please contact sales.'
        }
      }
    }
  },
  // Navigation translations
  nav: {
    home: 'Home',
    guide: 'Guide',
    membership: 'Membership',
    workbench: 'My Library',
  },

  // Layout component translations
  layout: {
    userProfile: 'Profile',
    settings: 'Settings',
    logout: 'Logout',
    login: 'Login',
  },

  // User Interface Elements
  ui: {
    login: 'Login',
    logout: 'Logout',
    profile: 'Profile',
    settings: 'Settings',
  },
  footer: {
    product: {
      title: 'Product',
      features: 'Features',
      usage: 'Usage Tips',
      userAgreement: 'User Agreement',
      privacy: 'Privacy Policy',
    },
    team: {
      title: 'Team',
      aboutUs: 'About Us',
      updates: 'Updates',
      contactUs: 'Contact Us',
      joinUs: 'Join Us',
    },
    social: {
      title: 'Social Media',
      one: 'Facebook',
      two: 'reddit',
      three: 'whatsapp'
    },
    copyright: 'All Rights Reserved',
  },
  TargetUserScenarios: {
    title: 'Target User Scenarios',
    tabs: {
      universityStudent: 'University Student',
      universityProfessor: 'University Professor',
      scientificResearcher: 'Scientific Researcher',
      corporateRDEngineer: 'Corporate R&D Engineer',
      lifelongLearner: 'Lifelong Learner',
    },
    content: {
      universityStudent: {
        title: 'Transforms a disorganized and weeks-long literature review process into a structured, efficient, and deep-thinking workflow.',
        description: [
          'Rapid Screening: Alex bulk-imports all 100 papers into the software. Using the <strong>“One-click AI Summary”</strong> feature, he quickly reads the AI-generated abstract for each paper. In just half a day, he narrows his core reading list down to 30 papers, saving a massive amount of time on initial screening.',
          'In-depth Reading: While reading a core paper, he encounters a paragraph full of complex formulas. He uses the <strong>“Ask questions about a specific selection”</strong> feature to ask AI, “Can you explain the physical meaning of this formula in simpler terms?” The AI provides a clear explanation, helping him overcome the barrier to understanding.',
          'Connecting Knowledge: After finishing each paper, he creates <strong>“Knowledge Cards”</strong> summarizing the core ideas, key figures, and his own insights. Finally, he arranges these cards on a visual whiteboard to map out the evolutionary relationships and a pro/con analysis of different algorithms, bringing his thesis structure to life.',
          'Standardized Writing: When writing his thesis, he uses the Word plugin to insert citations with a single click, automatically generating a <strong>“standard-format bibliography”</strong> that meets his university’s requirements.',
        ],
        scenario: {
          title: 'Scenario',
          text: 'He needs to read nearly a hundred English papers on the latest AI algorithms for his literature review and feels overwhelmed and inefficient.',
        },
        role: {
          title: 'Role',
          name: 'Alex',
          description: 'a Master\'s student in Computer Science, preparing his graduation thesis',
        },
        features: [
          'In-depth Reading',
          'Standardized Writing',
          'Rapid Screening',
          'Connecting Knowledge',
        ],
        imageUrl: '/images/university-student.png',
      },
      universityProfessor: {
        title: 'Dramatically increases her course prep efficiency and provides students with an innovative, guided, and interactive reading experience, shifting from one-way lectures to collaborative knowledge building.',
        description: [
          '<strong>Course Preparation:</strong> Professor Lee creates a shared library named "Digital Sociology Course." She adds dozens of core papers and book chapters required for the course.',
          '<strong>Creating Reading Guides:</strong> For each week\'s readings, she first reads the materials herself, highlighting key passages and adding notes as a guide. Then, she has the AI generate three critical thinking questions for each paper and appends them to her notes.',
          '<strong>Classroom Interaction:</strong> She shares the library with her students. They can read the annotated materials in advance and engage in group discussions based on the AI-generated questions, significantly improving the quality and interactivity of the class.',
        ],
        scenario: {
          title: 'Scenario',
          text: 'She needs to prepare course materials and create engaging learning experiences for her students while managing multiple classes.',
        },
        role: {
          title: 'Role',
          name: 'Professor Lee',
          description: 'a university professor teaching Digital Sociology',
        },
        features: [
          'Course Preparation',
          'Creating Reading Guides',
          'Classroom Interaction',
          'Collaborative Learning',
        ],
        imageUrl: '/images/university-professor.png',
      },
      scientificResearcher: {
        title: 'Elevates AI from a simple "Q&A tool" to an "intelligent partner" that sparks scientific inspiration, helping researchers efficiently discover innovation within a sea of information.',
        description: [
          '<strong>Discovering Innovation:</strong> Dr. Chen uses the <strong>"Ask questions about a screenshot of a figure"</strong> feature. She selects the data chart and asks the AI, "Beyond the author\'s conclusions, what other potential correlations or anomalies does this chart suggest?"',
          '<strong>Continuous Inquiry:</strong> The AI\'s response sparks an idea. She continues with <strong>"unlimited follow-up questions"</strong>: "Based on the anomaly you mentioned, what known signaling pathways might it be related to?" Through this continuous dialogue with the AI, she quickly identifies a previously overlooked research direction.',
          '<strong>Project Proposal:</strong> She organizes the insights from her AI-assisted brainstorming session into knowledge cards and links them to several other supporting papers, creating a logically sound and well-supported structured library that forms the solid foundation for her new research project proposal.',
        ],
        scenario: {
          title: 'Scenario',
          text: 'She needs to discover new research directions and develop innovative project proposals while analyzing complex scientific data.',
        },
        role: {
          title: 'Role',
          name: 'Dr. Chen',
          description: 'a researcher at a national research institute',
        },
        features: [
          'Discovering Innovation',
          'Continuous Inquiry',
          'Project Proposal',
          'AI-Assisted Analysis',
        ],
        imageUrl: '/images/scientific-researcher.png',
      },
      corporateRDEngineer: {
        title: 'Replaces the old, chaotic process of back-and-forth emails and confusing document versions with a centralized, traceable, and shared team knowledge base, significantly accelerating the team\'s technical decision-making process.',
        description: [
          '<strong>Team Knowledge Base:</strong> Zhang creates a shared team library and populates it with the technology\'s official white paper, core research papers, and several key technical blog posts.',
          '<strong>Asynchronous Collaboration:</strong> As team members read the documents, they add their questions and insights as notes directly on the PDFs. This allows everyone to see each other\'s thoughts, facilitating efficient asynchronous communication.',
          '<strong>Technical Comparison:</strong> During a team meeting, Zhang opens the software and filters for all member notes to guide the discussion. He also asks the AI on the spot: "Please create a table comparing this technology with our current solution, X, in terms of performance, cost, and community support." The AI quickly generates a clear and concise table.',
        ],
        scenario: {
          title: 'Scenario',
          text: 'He needs to coordinate team technical evaluations and make informed decisions while managing multiple stakeholders and complex technical information.',
        },
        role: {
          title: 'Role',
          name: 'Zhang',
          description: 'an R&D engineer at a tech company',
        },
        features: [
          'Team Knowledge Base',
          'Asynchronous Collaboration',
          'Technical Comparison',
          'Decision Making',
        ],
        imageUrl: '/images/corporate-r&d.png',
      },
      lifelongLearner: {
        title: 'The software acts as a patient, knowledgeable, and 24/7 "personal tutor," making profound scientific knowledge accessible and significantly lowering the barrier for lifelong learners to enter highly technical fields.',
        description: [
          '<strong>Lowering the Barrier to Entry:</strong> She finds a classic paper by Stephen Hawking but finds it full of jargon. She frequently uses the highlight-to-translate and targeted question features. For example, she highlights "event horizon" and asks, "Can you explain what this is using an analogy a layperson can understand?"',
          '<strong>Building a Knowledge System:</strong> Every time she grasps a new concept, she creates a knowledge card, summarizing it in her own words. Over time, these cards form her personal "black hole knowledge map."',
          '<strong>Exploratory Learning:</strong> After asking the AI about a concept, she follows up with, "Based on this, what should I learn about next?" The AI recommends related classic papers or concepts, guiding her through an exploratory learning journey.',
        ],
        scenario: {
          title: 'Scenario',
          text: 'As a software developer, she wants to understand complex astrophysics concepts from Stephen Hawking\'s papers in her spare time.',
        },
        role: {
          title: 'Role',
          name: 'Sarah',
          description: 'a lifelong learner with a passion for understanding complex scientific concepts',
        },
        features: [
          'Lowering the Barrier to Entry',
          'Building a Knowledge System',
          'Exploratory Learning',
          'Personal Tutoring',
        ],
        imageUrl: '/images/lifelong-learner.png',
      },
    }
  },
  DifferWithChatGPT: {
    title: 'How We Differ from ChatGPT',
    sections: {
      reading_retention: {
        title: 'Reading with Retention: Making Knowledge Stick',
        content: 'Unlike ChatGPT\'s conversational approach, odoc is designed specifically for academic literature. Our AI doesn\'t just answer questions—it helps you build lasting understanding through structured reading workflows, intelligent highlighting, and contextual note-taking that transforms how you absorb and retain complex information.'
      },
      professional_tool: {
        title: 'A Professional Academic Literature Tool',
        content: 'While ChatGPT excels at general conversations, odoc is purpose-built for researchers, students, and academics. We offer specialized features like citation management, literature organization, academic writing assistance, and research workflow optimization that ChatGPT simply cannot provide.'
      },
      structured_libraries: {
        title: 'Structured Libraries, Not Just Folders',
        content: 'ChatGPT stores conversations in simple lists. odoc creates intelligent knowledge networks. Our advanced organization system uses tags, categories, and AI-powered connections to help you discover relationships between papers, track research themes, and build comprehensive literature reviews effortlessly.'
      },
      immersive_experience: {
        title: 'An Immersive "What You See is What You Read" Experience',
        content: 'ChatGPT works through text exchanges. odoc provides a rich, visual reading environment where you can annotate directly on PDFs, see your highlights and notes in context, and experience seamless integration between reading, understanding, and knowledge creation—all in one unified interface.'
      },
      aiFeatures: {
        title: 'One-click AI Summary',
        content: 'For a long, unfamiliar paper, users can first generate a concise summary to quickly grasp its background, methods, findings, and conclusions.'
      },
      unlimitedQuestions: {
        title: 'Interactive Knowledge Flow: The Infinite Dialogue',
        content: 'This fosters a "Socratic" reading methodology. After the AI provides an initial answer, users can effortlessly pose deeper, multi-layered follow-up questions, creating a natural, continuous knowledge dialogue that transforms reading into a journey of critical discovery.'
      },
      specificSelection: {
        title: 'Precision Highlighting: Instant Micro-Comprehension',
        content: 'This represents the ultimate level of precision. When an ambiguous section or a complex phrase challenges your understanding, simply highlight the specific text. You can then directly interrogate the selection, obtaining immediate and highly focused answers tailored exactly to the highlighted segment.'
      }
    }
  },
  UserTestimonials: {
    title: 'User Testimonials',
    testimonials: {
      user1: {
        content: 'odoc is my favorite tool for reading and managing literature.\n\nThe AI-assisted reading feature is "exceptional". Reading tasks from my advisor that used to take me two or three days now only take half a day. I can ask the AI anything directly. It\'s a huge productivity booster! 👍\nI\'m always way ahead of my lab mates in every group meeting!',
        author: '@Researcher_TS8007A2',
        source: 'Source: User Community',
        avatar: 'user_1.svg'
      },
      user2: {
        content: 'I\'ve tried many similar products, but after all the searching, I always come back to odoc. It\'s just the smoothest experience for everything. It has helped me so much with reading and organizing papers. I love its literature management and AI-writing assistance, which makes my entire research and writing process faster and more efficient!',
        author: '@Researcher_TER4C0J2',
        source: 'Source: User Community',
        avatar: 'user_2.svg'
      },
      user3: {
        content: 'There are just too many great features to list!\n\nOrganizing with folders and tags makes finding papers incredibly fast.\n\nThe highlight-to-translate/full-text translation and AI-reader make reading papers a total breeze.\n\nThe note-taking and task management features make me super efficient.\nThanks, odoc!',
        author: '@Researcher_TSU1R562',
        source: 'Source: User Community',
        avatar: 'user_3.svg'
      },
      user4: {
        content: 'My favorite feature is the translation. As an international student, odoc\'s translation function has significantly eased the disadvantage of doing research in a non-native language. It has helped me maintain a first-class honors degree at a QS Top 15 university, one paper at a time. 🌟',
        author: '@momo',
        source: 'Source: User Community',
        avatar: 'user_4.svg'
      }
    }
  },
  OurPartners: {
    title: 'Our Partners',
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
}
