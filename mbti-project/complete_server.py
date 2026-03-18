#!/usr/bin/env python3
"""
MBTI Personality Test Server - Complete Version
支持用户注册、登录、真实题目测试、JWT认证
端口：8004
"""

import sqlite3
import hashlib
import json
import time
import jwt
import secrets
from datetime import datetime, timedelta, timezone
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
import os

# 配置
# 使用固定的SECRET_KEY，避免服务器重启后令牌失效
SECRET_KEY = "mbti_test_fixed_secret_key_2026_03_13_1234567890abcdef"  # 固定JWT密钥
DB_FILE = "mbti_test.db"
PORT = 8004

# 改进的MBTI题目（100道题 - 自然版）
MBTI_QUESTIONS = [
    {
        'id': 1,
        'question': "在聚会中，你通常：",
        'options': [
            {'text': "与很多人交流，包括陌生人", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "只与几个熟悉的人交谈", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 2,
        'question': "你更倾向于：",
        'options': [
            {'text': "通过实践学习", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "通过理论思考学习", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 3,
        'question': "做决定时，你更注重：",
        'options': [
            {'text': "逻辑和客观事实", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "情感和人际关系", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 4,
        'question': "你的生活风格更偏向：",
        'options': [
            {'text': "有计划和组织", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "灵活和随性", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 5,
        'question': "当你感到疲惫时，你更喜欢：",
        'options': [
            {'text': "与朋友一起放松", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "独自休息或阅读", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 6,
        'question': "你更相信：",
        'options': [
            {'text': "具体经验和事实", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "灵感和可能性", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 7,
        'question': "评价他人时，你更看重：",
        'options': [
            {'text': "他们的能力和成就", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "他们的感受和价值观", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 8,
        'question': "处理工作时，你倾向于：",
        'options': [
            {'text': "提前计划并按时完成", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "在压力下工作得更好", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 9,
        'question': "在新环境中，你通常：",
        'options': [
            {'text': "主动认识新朋友", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "先观察再参与", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 10,
        'question': "你更擅长：",
        'options': [
            {'text': "处理具体细节", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "看到整体大局", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 11,
        'question': "当有人批评你时，你更可能：",
        'options': [
            {'text': "分析批评是否有道理", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "感到受伤或不安", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 12,
        'question': "你的书桌通常是：",
        'options': [
            {'text': "整洁有序的", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "有些杂乱但你知道东西在哪", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 13,
        'question': "在社交场合，你通常：",
        'options': [
            {'text': "是谈话的中心", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "倾听多于发言", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 14,
        'question': "学习新东西时，你更喜欢：",
        'options': [
            {'text': "按部就班地学习", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "跳跃式地理解概念", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 15,
        'question': "做重要决定时，你更依赖：",
        'options': [
            {'text': "客观分析和数据", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "个人价值观和感受", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 16,
        'question': "你的假期计划通常是：",
        'options': [
            {'text': "详细安排好的", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "随性而定的", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 17,
        'question': "在团队中，你更可能：",
        'options': [
            {'text': "主动发言并带动气氛", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "倾听并思考后再表达", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 18,
        'question': "描述一个地方时，你更注重：",
        'options': [
            {'text': "具体的环境和细节", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "整体的氛围和感觉", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 19,
        'question': "当你需要放松时，你更倾向于：",
        'options': [
            {'text': "约朋友出去玩或参加社交活动", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "在家看书、看电影或独自休息", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 20,
        'question': "评价工作时，你更看重：",
        'options': [
            {'text': "效率和成果质量", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "团队合作和氛围", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 21,
        'question': "面对冲突时，你倾向于：",
        'options': [
            {'text': "就事论事分析问题", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "考虑各方感受", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 22,
        'question': "你的生活节奏更偏向：",
        'options': [
            {'text': "有计划有规律", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "灵活适应变化", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 23,
        'question': "学习新技能时，你更关注：",
        'options': [
            {'text': "技术的原理和逻辑", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "应用的价值和意义", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 24,
        'question': "学习新知识时，你更喜欢：",
        'options': [
            {'text': "按步骤实践掌握", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "理解背后的原理和概念", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 25,
        'question': "你的日常生活通常是：",
        'options': [
            {'text': "有计划、有规律的", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "灵活、随性的", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 26,
        'question': "处理多项任务时，你更擅长：",
        'options': [
            {'text': "按优先级顺序完成", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "多任务并行处理", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 27,
        'question': "你的工作环境通常是：",
        'options': [
            {'text': "整洁有序，物品归位", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "有些杂乱但你知道东西在哪", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 28,
        'question': "你的学习方式更偏向：",
        'options': [
            {'text': "小组讨论和互动学习", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "独立研究和自主学习", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 29,
        'question': "你的朋友圈通常是：",
        'options': [
            {'text': "广泛而多样，认识很多人", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "小而精，有几个知心好友", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 30,
        'question': "与人交流时，你更倾向于：",
        'options': [
            {'text': "谈论具体经历和事实", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "讨论想法和理论", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 31,
        'question': "在工作会议中，你通常：",
        'options': [
            {'text': "积极发言，表达自己的想法", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "先倾听，思考成熟后再发言", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 32,
        'question': "面对计划时，你更可能：",
        'options': [
            {'text': "严格遵守计划", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "根据情况调整", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 33,
        'question': "解决问题时，你更依赖：",
        'options': [
            {'text': "已知的有效方法", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "创新的解决方案", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 34,
        'question': "记忆信息时，你更容易记住：",
        'options': [
            {'text': "具体的细节和数据", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "整体的模式和联系", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 35,
        'question': "帮助他人时，你更可能：",
        'options': [
            {'text': "提供具体的解决方案", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "给予情感支持和理解", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 36,
        'question': "学习新技能时，你更喜欢：",
        'options': [
            {'text': "参加培训班与他人一起学习", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "自学或一对一指导", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 37,
        'question': "面对新事物时，你首先注意到：",
        'options': [
            {'text': "它的实际功能和外观", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "它的象征意义和潜力", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 38,
        'question': "面对批评时，你更可能：",
        'options': [
            {'text': "理性分析批评的合理性", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "感受批评带来的情绪影响", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 39,
        'question': "给予反馈时，你更注重：",
        'options': [
            {'text': "事实和准确性", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "方式和语气", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 40,
        'question': "你更关注：",
        'options': [
            {'text': "眼前的具体事实", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "未来的可能性", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 41,
        'question': "面对突然的变化时，你更可能：",
        'options': [
            {'text': "感到不安，希望恢复原计划", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "灵活适应，享受变化带来的新鲜感", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 42,
        'question': "在社交活动中，你通常：",
        'options': [
            {'text': "积极参与并享受热闹", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "更喜欢安静的小范围交流", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 43,
        'question': "你的价值观更偏向：",
        'options': [
            {'text': "公正和真理", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "和谐和同情", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 44,
        'question': "你的时间管理风格是：",
        'options': [
            {'text': "严格按时间表执行", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "大致安排，灵活调整", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 45,
        'question': "看待问题时，你倾向于：",
        'options': [
            {'text': "基于事实和经验分析", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "探索新的可能性和创意", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 46,
        'question': "面对选择时，你更可能：",
        'options': [
            {'text': "做出明确选择", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "保留多种可能性", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 47,
        'question': "认识新朋友时，你更可能：",
        'options': [
            {'text': "主动开启话题，热情交流", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "保持礼貌，等待对方先开口", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 48,
        'question': "周末你更喜欢：",
        'options': [
            {'text': "安排丰富的社交活动", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "享受安静的私人时间", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 49,
        'question': "处理工作任务时，你倾向于：",
        'options': [
            {'text': "提前规划，按时完成", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "灵活应对，在压力下效率更高", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 50,
        'question': "处理工作时，你更关注：",
        'options': [
            {'text': "任务完成的质量", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "团队合作的氛围", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 51,
        'question': "安排假期时，你倾向于：",
        'options': [
            {'text': "制定详细的行程计划", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "大致安排，留出自由发挥空间", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 52,
        'question': "你的能量主要来自：",
        'options': [
            {'text': "与外界互动和社交", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "独处和内省", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 53,
        'question': "处理人际关系时，你更倾向于：",
        'options': [
            {'text': "明确界限和原则", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "维护和谐和情感联系", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 54,
        'question': "做计划时，你更关注：",
        'options': [
            {'text': "具体的实施步骤和时间表", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "整体的目标和愿景", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 55,
        'question': "完成任务时，你倾向于：",
        'options': [
            {'text': "提前完成避免拖延", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "在截止日期前完成", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 56,
        'question': "选择礼物时，你更看重：",
        'options': [
            {'text': "实用性和具体功能", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "象征意义和情感价值", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 57,
        'question': "选择职业时，你更考虑：",
        'options': [
            {'text': "发展前景和薪资待遇", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "工作意义和团队氛围", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 58,
        'question': "描述事物时，你倾向于：",
        'options': [
            {'text': "使用具体细节和实例", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "使用比喻和象征", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 59,
        'question': "评价艺术作品时，你更注重：",
        'options': [
            {'text': "技巧和细节表现", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "情感表达和深层含义", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 60,
        'question': "做决定时，你更愿意：",
        'options': [
            {'text': "尽快做出明确决定", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "保持开放，收集更多信息", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 61,
        'question': "周末安排活动时，你倾向于：",
        'options': [
            {'text': "安排丰富的社交聚会", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "享受安静的私人时光", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 62,
        'question': "购物时，你更倾向于：",
        'options': [
            {'text': "有明确目标，快速完成", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "随意浏览，发现惊喜", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 63,
        'question': "当你需要思考问题时，你倾向于：",
        'options': [
            {'text': "通过与人讨论来理清思路", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "独自思考并整理想法", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 64,
        'question': "在团队项目中，你更擅长：",
        'options': [
            {'text': "协调沟通，带动团队氛围", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "深入研究，提供专业建议", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 65,
        'question': "学习知识时，你更重视：",
        'options': [
            {'text': "理论的严谨性", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "应用的关怀性", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 66,
        'question': "帮助他人时，你更倾向于：",
        'options': [
            {'text': "提供实际的解决方案", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "给予情感的支持", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 67,
        'question': "认识新朋友时，你通常：",
        'options': [
            {'text': "主动介绍自己并开启话题", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "等待对方先开口或观察", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 68,
        'question': "当你感到兴奋时，你更想：",
        'options': [
            {'text': "立刻与他人分享", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "先自己品味这份感受", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 69,
        'question': "旅行时，你更享受：",
        'options': [
            {'text': "体验具体的景点和活动", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "感受当地的文化和氛围", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 70,
        'question': "面对压力时，你更可能：",
        'options': [
            {'text': "找朋友倾诉寻求支持", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "独自处理情绪和问题", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 71,
        'question': "你的社交能量主要来自：",
        'options': [
            {'text': "与人互动和外部刺激", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "独处和内省的时间", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 72,
        'question': "组织活动时，你更注重：",
        'options': [
            {'text': "结构和流程", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "自由和创意", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 73,
        'question': "当你遇到有趣的事情时，你第一时间想：",
        'options': [
            {'text': "立刻分享给朋友或家人", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "先自己品味和思考", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 74,
        'question': "在工作环境中，你倾向于：",
        'options': [
            {'text': "开放式办公，便于交流", 'type': 'E', 'score_e': 1, 'score_i': 0},
            {'text': "独立空间，减少干扰", 'type': 'I', 'score_e': 0, 'score_i': 1},
        ]
    },
    {
        'id': 75,
        'question': "做决定时，你更重视：",
        'options': [
            {'text': "逻辑和客观标准", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "情感和人际关系", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 76,
        'question': "表达观点时，你更注重：",
        'options': [
            {'text': "逻辑的连贯性", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "情感的共鸣", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 77,
        'question': "给予他人反馈时，你更注重：",
        'options': [
            {'text': "事实的准确性和建设性", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "对方的感受和接受程度", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 78,
        'question': "做重要决定时，你更重视：",
        'options': [
            {'text': "逻辑分析和客观数据", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "个人感受和人际关系", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 79,
        'question': "制定规则时，你更注重：",
        'options': [
            {'text': "公平和一致性", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "灵活性和人性化", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 80,
        'question': "面对选择时，你更依赖：",
        'options': [
            {'text': "理性的分析", 'type': 'T', 'score_t': 1, 'score_f': 0},
            {'text': "内心的感受", 'type': 'F', 'score_t': 0, 'score_f': 1},
        ]
    },
    {
        'id': 81,
        'question': "处理信息时，你倾向于：",
        'options': [
            {'text': "快速分类做决定", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "持续收集新信息", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 82,
        'question': "学习时，你更喜欢：",
        'options': [
            {'text': "按步骤操作练习", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "理解背后的原理", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 83,
        'question': "你的决策过程更偏向：",
        'options': [
            {'text': "结论导向", 'type': 'J', 'score_j': 1, 'score_p': 0},
            {'text': "过程导向", 'type': 'P', 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 84,
        'question': "看待世界时，你更注重：",
        'options': [
            {'text': "现实和实际存在", 'type': 'S', 'score_s': 1, 'score_n': 0},
            {'text': "潜力和发展趋势", 'type': 'N', 'score_s': 0, 'score_n': 1},
        ]
    },
    {
        'id': 85,
        'question': "你更倾向于通过哪种方式学习新技能？",
        'options': [
            {'text': "参加培训课程或工作坊，与他人一起学习", 'type': 'E', 'score_e': 1, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "自己阅读资料或观看视频自学", 'type': 'I', 'score_e': 0, 'score_i': 1, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 86,
        'question': "对于一项新任务，你通常会：",
        'options': [
            {'text': "先制定详细的计划和时间表", 'type': 'J', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 1, 'score_p': 0},
            {'text': "边做边调整，保持灵活性", 'type': 'P', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 87,
        'question': "在团队讨论中，你更倾向于：",
        'options': [
            {'text': "提出大胆的创新想法", 'type': 'N', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 1, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "基于已有数据提出务实的方案", 'type': 'S', 'score_e': 0, 'score_i': 0, 'score_s': 1, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 88,
        'question': "当朋友向你倾诉烦恼时，你更倾向于：",
        'options': [
            {'text': "先给予情感上的支持和安慰", 'type': 'F', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 1, 'score_j': 0, 'score_p': 0},
            {'text': "帮他们分析问题并提供解决方案", 'type': 'T', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 1, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 89,
        'question': "你更喜欢哪种类型的工作环境？",
        'options': [
            {'text': "开放式的办公空间，方便随时交流", 'type': 'E', 'score_e': 1, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "安静的独立工作空间，减少干扰", 'type': 'I', 'score_e': 0, 'score_i': 1, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 90,
        'question': "在做重大人生决定时，你更看重：",
        'options': [
            {'text': "逻辑分析和客观事实", 'type': 'T', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 1, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "内心感受和价值观", 'type': 'F', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 1, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 91,
        'question': "面对假期安排，你通常会：",
        'options': [
            {'text': "提前规划好每天的行程", 'type': 'J', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 1, 'score_p': 0},
            {'text': "随性安排，到时再决定做什么", 'type': 'P', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 92,
        'question': "在阅读一本书时，你更享受：",
        'options': [
            {'text': "了解其中新颖的观点和理论", 'type': 'N', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 1, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "获取实用的知识和具体的方法", 'type': 'S', 'score_e': 0, 'score_i': 0, 'score_s': 1, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 93,
        'question': "在社交聚会中，你更可能是：",
        'options': [
            {'text': "主动与不认识的人攀谈交朋友", 'type': 'E', 'score_e': 1, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "与少数熟悉的朋友深入交流", 'type': 'I', 'score_e': 0, 'score_i': 1, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 94,
        'question': "对于未来，你更倾向于：",
        'options': [
            {'text': "畅想各种可能的未来场景", 'type': 'N', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 1, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "关注当下需要完成的具体事项", 'type': 'S', 'score_e': 0, 'score_i': 0, 'score_s': 1, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 95,
        'question': "当你需要评价一个人的表现时：",
        'options': [
            {'text': "会优先考虑他们的努力和付出的过程", 'type': 'F', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 1, 'score_j': 0, 'score_p': 0},
            {'text': "会以客观的结果和业绩来衡量", 'type': 'T', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 1, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 96,
        'question': "关于你的日常生活习惯：",
        'options': [
            {'text': "喜欢固定的作息和规律的生活节奏", 'type': 'J', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 1, 'score_p': 0},
            {'text': "喜欢每天都不太一样的灵活生活", 'type': 'P', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 1},
        ]
    },
    {
        'id': 97,
        'question': "你认为好的领导应该：",
        'options': [
            {'text': "关心团队成员的情感和福祉", 'type': 'F', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 1, 'score_j': 0, 'score_p': 0},
            {'text': "制定清晰的目标和高效的执行方案", 'type': 'T', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 1, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 98,
        'question': "当你收到一份复杂的说明书时：",
        'options': [
            {'text': "会仔细阅读每一个步骤", 'type': 'S', 'score_e': 0, 'score_i': 0, 'score_s': 1, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "会先了解大框架然后按自己的理解来操作", 'type': 'N', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 1, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 99,
        'question': "在周末休息时，你更愿意：",
        'options': [
            {'text': "参加社交活动或聚会", 'type': 'E', 'score_e': 1, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
            {'text': "在家享受安静的独处时光", 'type': 'I', 'score_e': 0, 'score_i': 1, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 0},
        ]
    },
    {
        'id': 100,
        'question': "对于项目的最终交付，你更倾向于：",
        'options': [
            {'text': "提前完成，留出缓冲时间检查", 'type': 'J', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 1, 'score_p': 0},
            {'text': "在截止日期前最后一刻完成", 'type': 'P', 'score_e': 0, 'score_i': 0, 'score_s': 0, 'score_n': 0, 'score_t': 0, 'score_f': 0, 'score_j': 0, 'score_p': 1},
        ]
    },
]

class MBTIServer(BaseHTTPRequestHandler):
    """MBTI测试服务器处理器"""
    
    def _init_db(self):
        """初始化数据库"""
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        # 创建用户表
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT UNIQUE NOT NULL,
                password_hash TEXT NOT NULL,
                nickname TEXT NOT NULL,
                gender TEXT NOT NULL,
                age INTEGER NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        ''')
        
        # 创建测试记录表
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS test_records (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                test_id TEXT NOT NULL,
                mbti_type TEXT,
                status TEXT DEFAULT 'in_progress',
                start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                end_time TIMESTAMP,
                FOREIGN KEY (user_id) REFERENCES users (id)
            )
        ''')
        
        # 创建答案详情表
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS answers (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                test_record_id INTEGER NOT NULL,
                question_id INTEGER NOT NULL,
                answer TEXT NOT NULL,
                score_e INTEGER DEFAULT 0,
                score_i INTEGER DEFAULT 0,
                score_s INTEGER DEFAULT 0,
                score_n INTEGER DEFAULT 0,
                score_t INTEGER DEFAULT 0,
                score_f INTEGER DEFAULT 0,
                score_j INTEGER DEFAULT 0,
                score_p INTEGER DEFAULT 0,
                FOREIGN KEY (test_record_id) REFERENCES test_records (id)
            )
        ''')
        
        conn.commit()
        conn.close()
    
    def _hash_password(self, password):
        """哈希密码"""
        return hashlib.sha256(password.encode()).hexdigest()
    
    def _verify_password(self, password, password_hash):
        """验证密码"""
        return self._hash_password(password) == password_hash
    
    def _generate_token(self, user_id, username):
        """生成JWT令牌"""
        payload = {
            'user_id': user_id,
            'username': username,
            'exp': datetime.now(timezone.utc) + timedelta(hours=24)
        }
        return jwt.encode(payload, SECRET_KEY, algorithm='HS256')
    
    def _verify_token(self, token):
        """验证JWT令牌"""
        try:
            payload = jwt.decode(token, SECRET_KEY, algorithms=['HS256'])
            return payload
        except:
            return None
    
    def _get_user_from_token(self):
        """从请求头获取用户信息"""
        auth_header = self.headers.get('Authorization')
        if not auth_header or not auth_header.startswith('Bearer '):
            return None
        
        token = auth_header.split(' ')[1]
        payload = self._verify_token(token)
        if not payload:
            return None
        
        return payload
    
    def _send_response(self, status_code, data=None, message=None):
        """发送JSON响应"""
        self.send_response(status_code)
        self.send_header('Content-Type', 'application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, PUT, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type, Authorization')
        self.end_headers()
        
        response = {}
        if message:
            response['message'] = message
        if data is not None:
            response.update(data)
        
        self.wfile.write(json.dumps(response, ensure_ascii=False).encode())
    
    def _send_error(self, status_code, message):
        """发送错误响应"""
        self._send_response(status_code, message=message)
    
    def _send_success(self, data=None, message="操作成功"):
        """发送成功响应"""
        self._send_response(200, data, message)
    
    def do_OPTIONS(self):
        """处理OPTIONS请求（CORS预检）"""
        self.send_response(200)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, PUT, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type, Authorization')
        self.end_headers()
    
    def do_GET(self):
        """处理GET请求"""
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        # 静态文件服务
        if path == '/' or path.endswith('.html') or path.endswith('.js') or path.endswith('.css'):
            self._serve_static_file(path)
            return
        
        # API路由
        if path == '/api/test/questions':
            self._get_questions()
        elif path == '/api/test/instructions':
            self._get_instructions()
        elif path == '/api/test/history':
            self._get_test_history()
        elif path == '/api/test/continue':
            self._get_continue_test()
        elif path == '/api/user/profile':
            self._get_user_profile()
        elif path == '/api/test/analysis':
            self._get_test_analysis()
        elif path == '/api/analytics/overview':
            self._get_analytics_overview()
        elif path == '/api/analytics/mbti-distribution':
            self._get_mbti_distribution()
        elif path == '/api/analytics/dimension-distribution':
            self._get_dimension_distribution()
        elif path == '/api/analytics/question-stats':
            self._get_question_stats()
        elif path == '/api/analytics/trends':
            self._get_analytics_trends()
        elif path == '/api/analytics/completion-rate':
            self._get_completion_rate()
        elif path == '/api/analytics/demographics':
            self._get_demographics()
        else:
            self._send_error(404, "接口不存在")
    
    def do_POST(self):
        """处理POST请求"""
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        content_length = int(self.headers.get('Content-Length', 0))
        if content_length > 0:
            body = self.rfile.read(content_length).decode('utf-8')
            try:
                data = json.loads(body)
            except:
                data = {}
        else:
            data = {}
        
        if path == '/api/register':
            self._register_user(data)
        elif path == '/api/login':
            self._login_user(data)
        elif path == '/api/test/start':
            self._start_test(data)
        elif path == '/api/test/submit':
            self._submit_test(data)
        else:
            self._send_error(404, "接口不存在")
    
    def do_PUT(self):
        """处理PUT请求"""
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        content_length = int(self.headers.get('Content-Length', 0))
        if content_length > 0:
            body = self.rfile.read(content_length).decode('utf-8')
            try:
                data = json.loads(body)
            except:
                data = {}
        else:
            data = {}
        
        if path == '/api/user/profile':
            self._update_user_profile(data)
        else:
            self._send_error(404, "接口不存在")
    
    def _serve_static_file(self, path):
        """提供静态文件"""
        if path == '/':
            file_path = 'index.html'
        else:
            file_path = path.lstrip('/')
        
        # 默认文件
        if not os.path.exists(file_path):
            if path == '/':
                file_path = 'login.html'
            elif path == '/login.html':
                file_path = 'login.html'
            elif path == '/test.html':
                file_path = 'test.html'
            else:
                self._send_error(404, "文件不存在")
                return
        
        try:
            with open(file_path, 'rb') as f:
                content = f.read()
            
            # 设置Content-Type
            if file_path.endswith('.html'):
                content_type = 'text/html'
            elif file_path.endswith('.js'):
                content_type = 'application/javascript'
            elif file_path.endswith('.css'):
                content_type = 'text/css'
            else:
                content_type = 'text/plain'
            
            self.send_response(200)
            self.send_header('Content-Type', content_type)
            self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
            self.send_header('Pragma', 'no-cache')
            self.send_header('Expires', '0')
            self.end_headers()
            self.wfile.write(content)
        except:
            self._send_error(500, "文件读取失败")
    
    def _register_user(self, data):
        """用户注册"""
        required_fields = ['username', 'password', 'nickname', 'gender', 'age']
        for field in required_fields:
            if field not in data or not str(data[field]).strip():
                self._send_error(400, f"缺少必填字段: {field}")
                return
        
        username = data['username'].strip()
        password = data['password']
        nickname = data['nickname'].strip()
        gender = data['gender'].strip()
        
        try:
            age = int(data['age'])
            if age < 1 or age > 120:
                raise ValueError
        except:
            self._send_error(400, "年龄必须是1-120之间的整数")
            return
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 检查用户名是否已存在
            cursor.execute("SELECT id FROM users WHERE username = ?", (username,))
            if cursor.fetchone():
                self._send_error(400, "用户名已存在")
                return
            
            # 插入新用户
            password_hash = self._hash_password(password)
            cursor.execute('''
                INSERT INTO users (username, password_hash, nickname, gender, age)
                VALUES (?, ?, ?, ?, ?)
            ''', (username, password_hash, nickname, gender, age))
            
            user_id = cursor.lastrowid
            conn.commit()
            
            # 生成令牌
            token = self._generate_token(user_id, username)
            
            self._send_success({
                'user_id': user_id,
                'username': username,
                'nickname': nickname,
                'token': token
            }, "注册成功")
        except Exception as e:
            self._send_error(500, f"注册失败: {str(e)}")
        finally:
            conn.close()
    
    def _login_user(self, data):
        """用户登录"""
        if 'username' not in data or 'password' not in data:
            self._send_error(400, "需要用户名和密码")
            return
        
        username = data['username'].strip()
        password = data['password']
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            cursor.execute('''
                SELECT id, username, password_hash, nickname, gender, age 
                FROM users WHERE username = ?
            ''', (username,))
            user = cursor.fetchone()
            
            if not user:
                self._send_error(401, "用户名或密码错误")
                return
            
            user_id, db_username, password_hash, nickname, gender, age = user
            
            if not self._verify_password(password, password_hash):
                self._send_error(401, "用户名或密码错误")
                return
            
            # 生成令牌
            token = self._generate_token(user_id, username)
            
            self._send_success({
                'user_id': user_id,
                'username': username,
                'nickname': nickname,
                'gender': gender,
                'age': age,
                'token': token
            }, "登录成功")
        except Exception as e:
            self._send_error(500, f"登录失败: {str(e)}")
        finally:
            conn.close()
    
    def _get_questions(self):
        """获取测试题目"""
        # 检查认证
        user = self._get_user_from_token()
        if not user:
            self._send_error(401, "需要登录")
            return
        
        # 返回题目（简化版，只返回必要信息）
        questions_for_client = []
        for q in MBTI_QUESTIONS:
            question_data = {
                'id': q['id'],
                'question': q['question'],
                'options': [{'text': opt['text'], 'type': opt['type']} for opt in q['options']]
            }
            questions_for_client.append(question_data)
        
        self._send_success({
            'questions': questions_for_client,
            'total': len(questions_for_client)
        })
    
    def _get_instructions(self):
        """获取测试说明"""
        instructions = {
            'title': 'MBTI性格测试',
            'description': '本测试基于经典的MBTI理论，通过16道题目帮助您了解自己的性格类型。',
            'steps': [
                '仔细阅读每道题目',
                '选择最符合您实际情况的选项',
                '请根据第一感觉选择，不要过多思考',
                '完成所有题目后提交，查看您的MBTI类型和分析'
            ],
            'duration': '约5-10分钟',
            'questions_count': len(MBTI_QUESTIONS)
        }
        self._send_success(instructions)
    
    def _get_test_history(self):
        """获取测试历史"""
        user = self._get_user_from_token()
        if not user:
            self._send_error(401, "需要登录")
            return
        
        user_id = user['user_id']
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            cursor.execute('''
                SELECT id, test_id, mbti_type, status, start_time, end_time
                FROM test_records 
                WHERE user_id = ? 
                ORDER BY start_time DESC
                LIMIT 10
            ''', (user_id,))
            
            records = []
            for row in cursor.fetchall():
                record_id, test_id, mbti_type, status, start_time, end_time = row
                records.append({
                    'id': record_id,
                    'test_id': test_id,
                    'mbti_type': mbti_type,
                    'status': status,
                    'start_time': start_time,
                    'end_time': end_time,
                    'duration': self._calculate_duration(start_time, end_time) if end_time else None
                })
            
            self._send_success({'records': records})
        except Exception as e:
            self._send_error(500, f"获取历史失败: {str(e)}")
        finally:
            conn.close()
    
    def _get_continue_test(self):
        """继续未完成的测试"""
        user = self._get_user_from_token()
        if not user:
            self._send_error(401, "需要登录")
            return
        
        user_id = user['user_id']
        
        # 检查是否指定了特定的 record_id
        record_id = self.query_params.get('record_id')
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            if record_id:
                # 获取特定测试记录
                cursor.execute('''
                    SELECT tr.id, tr.test_id, tr.start_time, tr.status
                    FROM test_records tr
                    WHERE tr.id = ? AND tr.user_id = ?
                ''', (record_id, user_id))
                
                record = cursor.fetchone()
                
                if record and record[3] == 'in_progress':
                    # 获取已答题目
                    cursor.execute('''
                        SELECT question_id, answer 
                        FROM answers 
                        WHERE test_record_id = ?
                        ORDER BY question_id
                    ''', (record_id,))
                    
                    answered = {row[0]: row[1] for row in cursor.fetchall()}
                    
                    self._send_success({
                        'test_id': record[1],
                        'record_id': record[0],
                        'start_time': record[2],
                        'answered': answered,
                        'can_continue': True
                    })
                else:
                    self._send_success({
                        'can_continue': False,
                        'message': '测试记录不存在或已完成'
                    })
            else:
                # 查找未完成的测试（默认行为）
                cursor.execute('''
                    SELECT tr.id, tr.test_id, tr.start_time
                    FROM test_records tr
                    WHERE tr.user_id = ? AND tr.status = 'in_progress'
                    ORDER BY tr.start_time DESC
                    LIMIT 1
                ''', (user_id,))
                
                record = cursor.fetchone()
                
                if record:
                    record_id, test_id, start_time = record
                    
                    # 获取已答题目
                    cursor.execute('''
                        SELECT question_id, answer 
                        FROM answers 
                        WHERE test_record_id = ?
                        ORDER BY question_id
                    ''', (record_id,))
                    
                    answered = {row[0]: row[1] for row in cursor.fetchall()}
                    
                    self._send_success({
                        'test_id': test_id,
                        'record_id': record_id,
                        'start_time': start_time,
                        'answered': answered,
                        'can_continue': True
                    })
                else:
                    self._send_success({
                        'can_continue': False,
                        'message': '没有未完成的测试'
                    })
                    'message': '没有未完成的测试'
                })
        except Exception as e:
            self._send_error(500, f"检查继续测试失败: {str(e)}")
        finally:
            conn.close()
    
    def _get_user_profile(self):
        """获取用户信息"""
        user = self._get_user_from_token()
        if not user:
            self._send_error(401, "需要登录")
            return
        
        user_id = user['user_id']
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            cursor.execute('''
                SELECT username, nickname, gender, age, created_at
                FROM users WHERE id = ?
            ''', (user_id,))
            
            user_data = cursor.fetchone()
            if user_data:
                username, nickname, gender, age, created_at = user_data
                
                # 获取测试统计
                cursor.execute('''
                    SELECT COUNT(*) as total_tests,
                           COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tests,
                           COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tests
                    FROM test_records WHERE user_id = ?
                ''', (user_id,))
                
                stats = cursor.fetchone()
                total_tests, completed_tests, in_progress_tests = stats
                
                self._send_success({
                    'username': username,
                    'nickname': nickname,
                    'gender': gender,
                    'age': age,
                    'created_at': created_at,
                    'stats': {
                        'total_tests': total_tests,
                        'completed_tests': completed_tests,
                        'in_progress_tests': in_progress_tests
                    }
                })
            else:
                self._send_error(404, "用户不存在")
        except Exception as e:
            self._send_error(500, f"获取用户信息失败: {str(e)}")
        finally:
            conn.close()
    
    def _update_user_profile(self, data):
        """更新用户信息"""
        user = self._get_user_from_token()
        if not user:
            self._send_error(401, "需要登录")
            return
        
        user_id = user['user_id']
        
        # 检查可更新字段
        updatable_fields = ['nickname', 'gender', 'age']
        update_data = {}
        
        for field in updatable_fields:
            if field in data:
                if field == 'age':
                    try:
                        age = int(data[field])
                        if age < 1 or age > 120:
                            raise ValueError
                        update_data[field] = age
                    except:
                        self._send_error(400, "年龄必须是1-120之间的整数")
                        return
                else:
                    value = str(data[field]).strip()
                    if value:
                        update_data[field] = value
        
        if not update_data:
            self._send_error(400, "没有可更新的字段")
            return
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 构建更新语句
            set_clause = ', '.join([f"{field} = ?" for field in update_data.keys()])
            values = list(update_data.values())
            values.append(user_id)
            
            cursor.execute(f'''
                UPDATE users SET {set_clause} WHERE id = ?
            ''', values)
            
            if cursor.rowcount > 0:
                conn.commit()
                self._send_success(message="用户信息更新成功")
            else:
                self._send_error(404, "用户不存在")
        except Exception as e:
            self._send_error(500, f"更新用户信息失败: {str(e)}")
        finally:
            conn.close()
    
    def _start_test(self, data):
        """开始测试"""
        user = self._get_user_from_token()
        if not user:
            self._send_error(401, "需要登录")
            return
        
        user_id = user['user_id']
        
        # 生成测试ID
        test_id = f"test_{int(time.time())}_{secrets.token_hex(4)}"
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 创建测试记录
            cursor.execute('''
                INSERT INTO test_records (user_id, test_id, status)
                VALUES (?, ?, 'in_progress')
            ''', (user_id, test_id))
            
            record_id = cursor.lastrowid
            conn.commit()
            
            self._send_success({
                'test_id': test_id,
                'record_id': record_id,
                'questions_count': len(MBTI_QUESTIONS),
                'message': '测试已开始'
            })
        except Exception as e:
            self._send_error(500, f"开始测试失败: {str(e)}")
        finally:
            conn.close()
    
    def _submit_test(self, data):
        """提交测试答案"""
        user = self._get_user_from_token()
        if not user:
            self._send_error(401, "需要登录")
            return
        
        user_id = user['user_id']
        
        if 'record_id' not in data or 'answers' not in data:
            self._send_error(400, "需要测试记录ID和答案")
            return
        
        record_id = data['record_id']
        answers = data['answers']
        
        # 验证答案格式
        if not isinstance(answers, dict):
            self._send_error(400, "答案格式不正确")
            return
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 验证测试记录属于当前用户
            cursor.execute('''
                SELECT id, status FROM test_records 
                WHERE id = ? AND user_id = ?
            ''', (record_id, user_id))
            
            record = cursor.fetchone()
            if not record:
                self._send_error(404, "测试记录不存在")
                return
            
            if record[1] == 'completed':
                self._send_error(400, "测试已完成，不能重复提交")
                return
            
            # 计算MBTI类型
            scores = {
                'E': 0, 'I': 0,
                'S': 0, 'N': 0,
                'T': 0, 'F': 0,
                'J': 0, 'P': 0
            }
            
            # 保存答案并计算分数
            for question_id_str, answer in answers.items():
                try:
                    question_id = int(question_id_str)
                except:
                    continue
                
                # 查找题目
                question = None
                for q in MBTI_QUESTIONS:
                    if q['id'] == question_id:
                        question = q
                        break
                
                if not question:
                    continue
                
                # 查找选项
                option = None
                for opt in question['options']:
                    if opt['text'] == answer:
                        option = opt
                        break
                
                if not option:
                    continue
                
                # 保存答案
                cursor.execute('''
                    INSERT INTO answers (
                        test_record_id, question_id, answer,
                        score_e, score_i, score_s, score_n,
                        score_t, score_f, score_j, score_p
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                ''', (
                    record_id, question_id, answer,
                    option.get('score_e', 0), option.get('score_i', 0),
                    option.get('score_s', 0), option.get('score_n', 0),
                    option.get('score_t', 0), option.get('score_f', 0),
                    option.get('score_j', 0), option.get('score_p', 0)
                ))
                
                # 累加分数
                for key in scores.keys():
                    score_key = f'score_{key.lower()}'
                    if score_key in option:
                        scores[key] += option[score_key]
            
            # 确定MBTI类型
            mbti_type = ''
            mbti_type += 'E' if scores['E'] > scores['I'] else 'I'
            mbti_type += 'S' if scores['S'] > scores['N'] else 'N'
            mbti_type += 'T' if scores['T'] > scores['F'] else 'F'
            mbti_type += 'J' if scores['J'] > scores['P'] else 'P'
            
            # 更新测试记录
            cursor.execute('''
                UPDATE test_records 
                SET mbti_type = ?, status = 'completed', end_time = CURRENT_TIMESTAMP
                WHERE id = ?
            ''', (mbti_type, record_id))
            
            conn.commit()
            
            # 获取分析报告
            analysis = self._generate_analysis(mbti_type, scores)
            
            self._send_success({
                'mbti_type': mbti_type,
                'scores': scores,
                'analysis': analysis,
                'record_id': record_id,
                'message': '测试提交成功'
            })
        except Exception as e:
            self._send_error(500, f"提交测试失败: {str(e)}")
        finally:
            conn.close()
    
    def _get_test_analysis(self):
        """获取测试分析"""
        user = self._get_user_from_token()
        if not user:
            self._send_error(401, "需要登录")
            return
        
        parsed_path = urlparse(self.path)
        query_params = parse_qs(parsed_path.query)
        
        record_id = query_params.get('record_id', [None])[0]
        if not record_id:
            self._send_error(400, "需要测试记录ID")
            return
        
        try:
            record_id = int(record_id)
        except:
            self._send_error(400, "测试记录ID格式不正确")
            return
        
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 验证测试记录属于当前用户
            cursor.execute('''
                SELECT tr.mbti_type, tr.start_time, tr.end_time
                FROM test_records tr
                WHERE tr.id = ? AND tr.user_id = ? AND tr.status = 'completed'
            ''', (record_id, user['user_id']))
            
            record = cursor.fetchone()
            if not record:
                self._send_error(404, "测试记录不存在或未完成")
                return
            
            mbti_type, start_time, end_time = record
            
            # 获取分数详情
            cursor.execute('''
                SELECT 
                    SUM(score_e) as total_e,
                    SUM(score_i) as total_i,
                    SUM(score_s) as total_s,
                    SUM(score_n) as total_n,
                    SUM(score_t) as total_t,
                    SUM(score_f) as total_f,
                    SUM(score_j) as total_j,
                    SUM(score_p) as total_p
                FROM answers WHERE test_record_id = ?
            ''', (record_id,))
            
            score_row = cursor.fetchone()
            scores = {
                'E': score_row[0] or 0, 'I': score_row[1] or 0,
                'S': score_row[2] or 0, 'N': score_row[3] or 0,
                'T': score_row[4] or 0, 'F': score_row[5] or 0,
                'J': score_row[6] or 0, 'P': score_row[7] or 0
            }
            
            # 生成分析报告
            analysis = self._generate_analysis(mbti_type, scores)
            
            self._send_success({
                'mbti_type': mbti_type,
                'scores': scores,
                'analysis': analysis,
                'start_time': start_time,
                'end_time': end_time,
                'duration': self._calculate_duration(start_time, end_time)
            })
        except Exception as e:
            self._send_error(500, f"获取分析失败: {str(e)}")
        finally:
            conn.close()
    
    dimension_descriptions = {
        'E': "外向(E): 你在社交互动中获取能量, 享受与人交流和合作。你喜欢热闹的环境, 善于表达想法, 在团队中通常是活跃的参与者。你的社交圈广泛, 从各种人际交往中汲取灵感和动力。在工作和生活中, 你倾向于通过讨论和协作来解决问题, 面对面的交流比独自思考更能激发你的创造力。",
        'I': "内向(I): 你在独处和安静的环境中获取能量, 深度思考是你最自然的状态。你喜欢深入的、一对一的交流而非大型社交活动, 拥有少数但非常亲密的友谊。你的内心世界丰富而活跃, 独处对你来说不是寂寞而是充电。在工作和生活中, 你倾向于先思考再发言, 注重质量和深度胜过数量。",
        'S': "实感(S): 你倾向于关注具体的事实、细节和当下的实际情况。你注重实际经验, 相信眼见为实, 善于观察和记忆具体的信息。你做事脚踏实地, 偏好已知和经过验证的方法。在解决问题时, 你从具体细节出发, 逐步构建解决方案, 对现实世界的运作方式有着敏锐的感知。",
        'N': "直觉(N): 你倾向于关注事物的可能性、模式和潜在含义, 而非表面的细节。你享受想象未来的各种可能性, 善于从现有信息中推导出新的理论和洞察。你喜欢抽象思维和理论探讨, 对创新和变革充满热情。在解决问题时, 你习惯从宏观和长远的角度来思考, 能够看到别人忽略的联系和机会。",
        'T': "思考(T): 你在做决定时主要依赖逻辑分析和客观事实。你追求公平和一致性, 善于用理性的方式评估各种选择的利弊。你不会被个人情感过度影响判断, 这使得你的决策通常冷静而客观。在处理问题时, 你优先考虑效率和正确性, 相信理性分析能够带来最好的结果。",
        'F': "情感(F): 你在做决定时主要考虑价值观、人际关系和对他人的影响。你有出色的同理心, 能够站在他人的角度思考问题。你追求和谐的人际关系, 倾向于选择对大多数人都有利的方案。在处理问题时, 你优先考虑人们的感受和需求, 相信善良和理解能够带来最好的结果。",
        'J': "判断(J): 你偏好有计划、有条理的生活方式。你喜欢提前规划和安排, 享受完成计划带来的成就感。你做事有明确的目标和时间表, 讨厌意外和混乱。你倾向于快速做出决定, 保持事物的有序进行。在工作方面, 你擅长制定计划、分配资源和把控进度。",
        'P': "感知(P): 你偏好灵活、开放的生活方式。你喜欢保持选择的开放性, 享受在最后时刻做出决定带来的灵活性。你适应能力强, 善于在变化中找到新的机会。你倾向于收集更多信息再做出判断, 享受探索各种可能性的过程。在工作方面, 你擅长应对突发情况和多任务处理。",
    }

    dimension_descriptions = {
        'E': "外向(E): 你在社交互动中获取能量, 享受与人交流和合作。你喜欢热闹的环境, 善于表达想法, 在团队中通常是活跃的参与者。你的社交圈广泛, 从各种人际交往中汲取灵感和动力。在工作和生活中, 你倾向于通过讨论和协作来解决问题, 面对面的交流比独自思考更能激发你的创造力。",
        'I': "内向(I): 你在独处和安静的环境中获取能量, 深度思考是你最自然的状态。你喜欢深入的、一对一的交流而非大型社交活动, 拥有少数但非常亲密的友谊。你的内心世界丰富而活跃, 独处对你来说不是寂寞而是充电。在工作和生活中, 你倾向于先思考再发言, 注重质量和深度胜过数量。",
        'S': "实感(S): 你倾向于关注具体的事实、细节和当下的实际情况。你注重实际经验, 相信眼见为实, 善于观察和记忆具体的信息。你做事脚踏实地, 偏好已知和经过验证的方法。在解决问题时, 你从具体细节出发, 逐步构建解决方案, 对现实世界的运作方式有着敏锐的感知。",
        'N': "直觉(N): 你倾向于关注事物的可能性、模式和潜在含义, 而非表面的细节。你享受想象未来的各种可能性, 善于从现有信息中推导出新的理论和洞察。你喜欢抽象思维和理论探讨, 对创新和变革充满热情。在解决问题时, 你习惯从宏观和长远的角度来思考, 能够看到别人忽略的联系和机会。",
        'T': "思考(T): 你在做决定时主要依赖逻辑分析和客观事实。你追求公平和一致性, 善于用理性的方式评估各种选择的利弊。你不会被个人情感过度影响判断, 这使得你的决策通常冷静而客观。在处理问题时, 你优先考虑效率和正确性, 相信理性分析能够带来最好的结果。",
        'F': "情感(F): 你在做决定时主要考虑价值观、人际关系和对他人的影响。你有出色的同理心, 能够站在他人的角度思考问题。你追求和谐的人际关系, 倾向于选择对大多数人都有利的方案。在处理问题时, 你优先考虑人们的感受和需求, 相信善良和理解能够带来最好的结果。",
        'J': "判断(J): 你偏好有计划、有条理的生活方式。你喜欢提前规划和安排, 享受完成计划带来的成就感。你做事有明确的目标和时间表, 讨厌意外和混乱。你倾向于快速做出决定, 保持事物的有序进行。在工作方面, 你擅长制定计划、分配资源和把控进度。",
        'P': "感知(P): 你偏好灵活、开放的生活方式。你喜欢保持选择的开放性, 享受在最后时刻做出决定带来的灵活性。你适应能力强, 善于在变化中找到新的机会。你倾向于收集更多信息再做出判断, 享受探索各种可能性的过程。在工作方面, 你擅长应对突发情况和多任务处理。",
    }

    def _generate_analysis(self, mbti_type, scores):
        type_descriptions = {
            'ISTJ': "检查员型(ISTJ) - 作为MBTI中最靠谱的类型之一, ISTJ是社会的基石。你安静、严肃, 通过全面性和可靠性获得成功。你注重实际, 有高度的责任感, 做事脚踏实地, 从不半途而废。在你的世界里, 秩序和规则不是束缚, 而是让一切高效运转的保障。你拥有惊人的记忆力和对细节的把控能力, 能够回忆起别人早已遗忘的重要信息。你倾向于用事实和逻辑来处理问题, 不喜欢凭空猜测。你的承诺一旦给出, 就如同刻在石头上一样坚定 - 你是最值得信赖的朋友和同事。在工作中, 你擅长建立和维护系统化的流程, 确保每个环节都不出差错。你的思维方式偏向线性和结构化, 善于将复杂问题拆解为可执行的步骤。你的行为模式体现着言出必行的信念, 说到做到是你的人生信条。你的核心优势包括极高的责任感和可靠性, 出色的组织和管理能力, 强大的记忆力和对细节的关注, 以及在压力下依然保持冷静和理性的能力。你天生适合需要精确和一致性的工作。你的潜在盲点在于可能过于固执于既定规则, 对变化和新事物的接受速度较慢, 有时过于严苛地要求自己和他人都达到完美标准, 在表达情感方面可能显得含蓄甚至冷淡, 容易忽略他人的情感需求。成长方向建议你尝试有意识地接受变化, 给自己留出灵活性空间; 练习表达情感, 让关心你的人感受到你的温暖; 在坚持原则的同时, 学会倾听不同的声音和观点; 适当放松对自己的要求, 认识到足够好有时比完美更实际。",
            'ISFJ': "照顾者型(ISFJ) - ISFJ是最温暖、最体贴的人格类型之一。你安静、友好、有责任心、有良知, 是你身边人最坚实的后盾。你坚定地致力于完成你的义务, 无论是对家人、朋友还是工作团队。你拥有超强的观察力, 能够敏锐地察觉到他人的需求和情绪变化, 并在别人开口之前就主动提供帮助。你的记忆力惊人, 尤其擅长记住对别人来说重要的事情 - 生日、偏好、说过的话 - 这些细节让你成为最贴心的存在。你忠于传统和价值观, 倾向于用实际行动而非言语来表达关爱。你做事认真细致, 追求完美, 在默默付出的过程中获得巨大的满足感。你不喜欢成为注意力的中心, 更愿意在幕后确保一切井然有序。你的善良不是软弱, 而是一种深沉的力量, 能够在困难时刻为身边的人提供稳固的支持。你的核心优势包括出色的同理心和关怀能力, 细致入微的观察力, 强大的记忆力和对细节的关注, 高度的责任感和忠诚度, 以及营造和谐氛围的天赋。你是团队中最可靠的成员, 也是最温暖的朋友。你的潜在盲点在于可能过度牺牲自己的需求来满足他人, 难以拒绝别人的请求导致自己负担过重, 对变化和不确定性的适应能力较弱, 倾向于压抑自己的负面情绪, 在需要表达不同意见时可能过于犹豫。成长方向建议你学会说不, 认识到自我照顾不是自私; 勇敢面对变化, 尝试走出舒适区; 定期审视自己的情绪和需求, 不要总是把自己放在最后; 当有不同意见时, 学会以温和但坚定的方式表达出来。",
            'INFJ': "提倡者型(INFJ) - INFJ是MBTI中最稀有的类型之一, 拥有独特的深度和洞察力。你寻求思想、关系和物质之间的意义和联系, 对事物的本质有着与生俱来的敏锐感知。你有惊人的直觉, 常常能在一瞬间看透事情的本质和人的内心世界。你对如何更好地服务大众有着清晰的远景, 心中燃烧着改变世界的热情。你是理想主义者与行动派的完美结合 - 不仅有崇高的理想, 还愿意付诸实践去实现它们。你重视深度胜过广度, 宁可拥有几个知心好友, 也不愿参与表面的社交活动。你有丰富的内心世界, 充满想象力, 对文学、艺术和哲学有着天然的兴趣。你对不公正的事情有着强烈的感受, 会为信念挺身而出。你善于用文字和语言来表达复杂的思想和情感, 你的写作和表达能力常常令人惊叹。你的核心优势包括卓越的洞察力和直觉能力, 深刻的同理心和理解他人情感的能力, 强大的创造力和想象力, 坚定的价值观和理想主义精神, 以及能够激励和感染他人的沟通能力。你的潜在盲点在于可能对自己和他人期望过高, 容易感到失望和理想破灭; 对冲突非常敏感, 倾向于回避直接的对抗; 可能过于封闭内心, 不轻易让人接近; 在追求完美的过程中容易忽略现实条件; 有时会忽略自己的生理需求。成长方向建议你接受世界和人的不完美, 学会在不完美中找到美好; 面对冲突时不要逃避, 它是解决问题的必要途径; 学会向信任的人敞开心扉; 在追求理想的同时关注当下的实际需求。",
            'INTJ': "建筑师型(INTJ) - INTJ是最独立、最有战略眼光的人格类型之一。你在实现自己的想法和达成目标时有创新的想法和非凡的动力。你的大脑就像一台高速运转的超级计算机, 不断分析、规划和优化。你有天生的系统思维, 能够看到别人看不到的宏大图景和深层联系。你高度独立, 不依赖外界的认可和指导, 自信地按照自己的判断行事。你对知识和能力有着永无止境的渴望, 享受攻克智力难题的快感。你做决策时冷静理性, 不会被情感左右, 这让你在面对复杂问题时能够做出最有效的选择。你有天生的领导才能, 但更倾向于通过专业能力和远见卓识来影响他人, 而非传统的管理方式。你追求效率和卓越, 对低效和愚蠢的事情缺乏耐心。你的内心世界丰富而有序, 充满了各种有趣的想法和计划。你的核心优势包括卓越的战略规划能力, 出色的逻辑分析能力, 强大的独立思考和创新精神, 高度的自律和执行力, 以及快速学习和掌握复杂系统的能力。你天生就是解决难题的高手。你的潜在盲点在于可能过于理性而忽略情感因素, 在人际交往中显得冷漠或傲慢; 对标准较低的人缺乏耐心和包容; 可能过于关注未来目标而忽略当下的感受; 在表达情感方面存在困难; 可能低估他人的情感需求。成长方向建议你有意识地关注他人的情感需求, 培养同理心; 学会欣赏不同人的不同优势, 而非仅以能力标准衡量; 在追求远大目标的同时享受当下的过程; 练习表达感受, 让亲近的人更了解你。",
            'ISTP': "鉴赏家型(ISTP) - ISTP是MBTI中最灵活、最实际的人格类型之一。你灵活、忍耐力强, 是个安静的观察者。你在日常中看似随性, 但一旦有问题发生, 就会立刻行动, 找到最实用的解决方法。你有着与生俱来的机械天赋, 喜欢理解事物是如何运作的, 擅长将复杂系统分解为基本组成部分。你热爱动手实践, 相信实践出真知, 对你来说, 纸上谈兵毫无意义。你冷静、理性, 在紧急情况下反而更加沉着, 是危机时刻最能依靠的人。你注重效率, 讨厌不必要的规则和繁文缛节, 倾向于用最直接的方式解决问题。你享受冒险和刺激, 喜欢挑战自己的极限。你有强烈的个人空间需求, 重视自由和独立胜过一切。你的好奇心驱使你不断尝试新事物, 学习新技能。你的核心优势包括出色的动手能力和技术天赋, 冷静沉着的危机处理能力, 灵活适应各种情况的能力, 敏锐的观察力和分析能力, 以及将理论转化为实际行动的高效执行力。你是最务实的问题解决者。你的潜在盲点在于可能过于享受当下而忽略长远规划, 对承诺和义务感到束缚; 在表达情感方面比较困难; 可能过于独立而难以建立深层的人际关系; 有时会显得冷漠或疏离; 可能突然改变计划, 让身边的人感到不安。成长方向建议你学会为未来做一些规划和准备; 有意识地表达情感, 让身边的人知道你关心他们; 在追求自由的同时, 认识到承诺和责任的价值; 培养持久的人际关系, 享受深度连接带来的满足感。",
            'ISFP': "探险家型(ISFP) - ISFP是MBTI中最具艺术气质和温暖灵魂的人格类型之一。你安静、友好、敏感、和善, 拥有独特的审美眼光和丰富的内心世界。你享受当下, 对生活有着敏锐的感知力 - 一朵花的绽放、一缕阳光的角度、一首曲子的旋律, 这些微小的美好在你的眼中都有非凡的意义。你忠于自己的价值观, 虽然不一定说出来, 但会用行动坚定地守护你认为重要的事物。你有空间感和审美天赋, 很多人在音乐、绘画、设计等领域有着出色的才华。你随和自在, 不喜欢冲突和对抗, 倾向于以温和的方式化解矛盾。你享受有自己的空间, 按照自己的节奏工作和生活。你对人和事物有着深刻的感知力, 能够感受到别人忽略的细微之处。你的善良和真诚是发自内心的, 不需要刻意表现。你的核心优势包括出色的审美能力和艺术创造力, 强烈的同理心和他人的关怀, 活在当下享受生活的能力, 灵活适应不同环境的天赋, 以及通过行动而非言语表达爱的温暖品质。你的潜在盲点在于可能过于回避冲突, 导致自己的需求被忽略; 缺乏长期规划和目标设定; 有时过于敏感, 容易受到批评的伤害; 可能过于随性而显得不够可靠; 在面对需要果断决策的场合可能犹豫不决。成长方向建议你学会在必要时表达自己的需求和观点, 即使这意味着可能产生冲突; 制定一些长期目标, 给自己的生活方向感; 培养接受建设性批评的能力; 在保持随性的同时, 适度增加生活的结构和规划。",
            'INFP': "调解员型(INFP) - INFP是MBTI中最富诗意和理想主义色彩的人格类型。你理想主义, 对于自己的价值观和自己觉得重要的人非常忠诚。你内心深处燃烧着让世界变得更好的渴望, 希望外部的生活和内心的价值观达到完美的统一。你有丰富的想象力和创造力, 脑海中常常浮现出各种美好的愿景。你善于用文字、艺术和其他创造性方式来表达内心深处的情感和思想。你是深度的思考者, 不满足于表面答案, 总是追寻事物背后的意义。你对人有着真诚的兴趣, 喜欢了解每个人的故事和内心世界。你有强烈的道德感, 当看到不公正的事情时会挺身而出。你安静但不怯弱, 内敛但充满力量。你的情感丰富而深沉, 虽然不一定外露, 但对亲近的人来说, 你的爱和关怀是深沉而持久的。你的核心优势包括卓越的创造力和想象力, 深厚的同理心和理解能力, 坚定的价值观和道德准则, 出色的写作和表达能力, 以及激励他人追求理想的感染力。你的潜在盲点在于可能过于理想化, 对现实的失望容易导致消极情绪; 在实现理想时缺乏行动力和执行力; 可能过于沉浸在自己的内心世界而忽略外部环境; 对自己和对他人的期望有时不够现实; 在面对日常琐事和行政事务时感到厌烦。成长方向建议你将理想分解为具体可执行的小步骤, 逐步实现; 培养日常的行动习惯, 不只是停留在思考阶段; 接受不完美是生活的一部分; 在追求精神世界的同时也要照顾好现实生活的需求。",
            'INTP': "逻辑学家型(INTP) - INTP是MBTI中最具逻辑性和探索精神的人格类型。你对于自己感兴趣的任何事物都寻求找到合理的解释, 你的大脑是一个永不停歇的思想实验室。你喜欢理论和抽象的事物, 享受在理念的世界中自由探索。对你来说, 一个精妙的想法比任何物质奖励都更有价值。你有天生的分析能力, 能够迅速识别逻辑漏洞和不一致之处。你不满足于表层的答案, 总是追问为什么和如果会怎样。你有广泛的兴趣, 能够在多个领域之间建立意想不到的联系。你喜欢独立思考, 不随波逐流, 对权威和传统持有健康的质疑态度。你有出色的抽象思维能力, 能够理解和构建复杂的理论体系。你善于发现模式和规律, 在科学、数学、哲学和技术领域有着天然的优势。虽然你外表安静, 但内心世界异常活跃, 充满了各种有趣的想法和理论。你的核心优势包括卓越的逻辑分析能力和批判性思维, 强大的抽象思维和理论建构能力, 快速学习和掌握复杂概念的能力, 创造性的问题解决方式, 以及对知识的纯粹热爱和追求。你的潜在盲点在于可能过于沉浸在思考中而忽略了实际行动; 在社交场合可能显得疏离或不善表达; 对重复性和日常性的工作缺乏耐心; 可能过于追求完美理论而忽略实际可行性; 在表达情感方面存在困难, 让人觉得难以接近。成长方向建议你将好的想法付诸实践, 不要只停留在理论层面; 有意识地培养社交技能, 与志同道合的人建立联系; 接受足够好的解决方案, 不必追求每个决定都完美无缺; 练习表达情感, 让关心你的人了解你的内心。",
            'ESTP': "企业家型(ESTP) - ESTP是MBTI中最充满活力和冒险精神的人格类型。你灵活、忍耐力强, 实际, 注重结果。你觉得理论和抽象的解释非常无趣, 对你来说, 只有看得见摸得着的行动才有意义。你喜欢积极地采取行动解决问题, 是典型的先做再说类型。你有天生的社交魅力, 能够迅速与各种各样的人建立联系。你观察力敏锐, 能够快速读取环境和人的信息, 并做出即时反应。你享受冒险和挑战, 生活对你来说就是一场精彩的冒险。你有出色的应变能力, 在突发情况下能够快速调整策略。你直率坦诚, 不喜欢拐弯抹角, 这让你在某些人看来可能过于直接, 但也赢得了许多人的尊重和信任。你注重当下, 善于把握眼前的机会。你有天生的商业头脑和谈判能力, 能够在竞争激烈的环境中脱颖而出。你的核心优势包括出色的应变和危机处理能力, 强大的观察力和快速判断力, 天生的社交魅力和影响力, 高效的实际执行能力, 以及在压力下保持冷静和果断的能力。你的潜在盲点在于可能过于追求刺激而忽略长期后果; 对规则和流程缺乏耐心, 可能铤而走险; 在需要深入思考和耐心的工作中容易感到厌倦; 可能过于关注当下而忽略未来规划; 在处理复杂情感问题时显得笨拙。成长方向建议你在行动之前花一点时间思考可能的后果; 培养耐心, 学会在需要坚持的任务中保持专注; 为自己制定一些中长期目标; 学会欣赏理论思考和计划的重要性。",
            'ESFP': "表演者型(ESFP) - ESFP是MBTI中最具感染力和热情的人格类型。你外向、友好、接受力强, 热爱生活、热爱人类、热爱一切美好的事物。你的存在就像一束阳光, 走到哪里都能照亮周围。你有天生的表演天赋和娱乐精神, 能够让任何场合变得生动有趣。你对美的感受力极强, 无论是视觉、听觉还是味觉, 你都追求最优质的体验。你注重当下, 善于发现和创造快乐。你是天生的社交达人, 能够轻松地与任何人建立联系。你有强烈的同理心, 能够感知到他人的情绪变化, 并做出恰当的回应。你慷慨大方, 乐于分享自己的时间和资源。你善于实际操作, 喜欢用行动来表达关心和爱意。你有出色的即兴表现能力, 即使在意外情况下也能表现得游刃有余。你的乐观和热情是极具感染力的, 身边的人都喜欢和你在一起。你的核心优势包括出色的社交能力和人际魅力, 强大的即兴应变和适应能力, 对美的敏锐感知和审美能力, 激励他人和创造快乐的天赋, 以及将想法迅速转化为实际行动的执行力。你的潜在盲点在于可能过于追求享乐而忽视责任和义务; 在面对需要深度思考和长期规划的任务时容易感到困难; 可能过于关注他人的看法而忽略自己的真实需求; 有时缺乏耐心处理繁琐的细节工作; 在冲突面前倾向于回避而非正面解决。成长方向建议你在享受当下的同时, 也为未来做一些储备和规划; 培养深度思考的习惯, 不要只停留在表面; 学会独处和自我反思, 这将帮助你更深地了解自己; 在人际交往中, 也要关注自己的需求, 不要总是迎合他人。",
            'ENFP': "竞选者型(ENFP) - ENFP是MBTI中最具热情和创造力的人格类型。你热情洋溢、富有想象力, 认为人生充满了无限的可能性。你的大脑像一台创意永动机, 不断产生新的想法、新的计划、新的可能性。你能很快地将事物和信息联系起来, 然后自信地根据自己的判断解决问题。你有天生的感染力, 能够用你的热情和乐观激励周围的人。你善于发现每个人的潜力和闪光点, 这让身边的人在你面前感到被认可和鼓舞。你重视真实和自由, 不喜欢虚伪和束缚。你对世界充满好奇, 总是想要了解更多的人、更多的事、更多的领域。你有出色的语言表达能力, 无论是写作还是说话, 都能够打动人心。你善于创新和变革, 是推动进步的重要力量。你的温暖和真诚让你在任何社交场合都如鱼得水。你的核心优势包括卓越的创造力和想象力, 强大的感染力和激励他人的能力, 出色的沟通和表达能力, 快速的思维敏捷性和适应能力, 以及发现可能性和机会的敏锐眼光。你的潜在盲点在于可能同时追求太多可能性, 导致精力分散、难以专注完成一件事; 在实际执行阶段容易失去兴趣; 对日常事务和细节工作缺乏耐心; 可能过于理想化, 对现实的困难估计不足; 有时过于热情而给人不够成熟的印象。成长方向建议你选择最重要的目标, 学会专注和坚持到底; 培养完成项目的纪律性, 不要只停留在创意阶段; 认识到日常事务是创意实现的基础; 学会区分热情和深度承诺的区别。",
            'ENTP': "辩论家型(ENTP) - ENTP是MBTI中最聪明、最有辩论精神的人格类型。你反应快、睿智, 有激励别人的能力, 警觉性强、直言不讳。你的大脑像一台永不停歇的创意引擎, 不断挑战现有观念, 提出新的假说和可能性。你享受智力上的交锋, 不是因为好胜, 而是因为在辩论中能够发现更深层的真理。你有天生的系统分析能力, 能够快速理解复杂的概念和框架。你善于从不同角度看问题, 这让你的解决方案常常出人意料但极为有效。你有强烈的好奇心和求知欲, 对任何话题都感兴趣, 这使得你拥有广博的知识面。你不畏惧挑战权威和传统, 敢于提出不同的声音。你有出色的语言表达能力和幽默感, 能够在社交场合中轻松驾驭各种话题。你善于发现系统中的弱点和改进空间, 是天生的创新者和问题解决者。你的核心优势包括卓越的辩论和分析能力, 强大的创新和系统思维能力, 快速学习和适应新领域的能力, 出色的语言表达和沟通能力, 以及发现问题和机会的敏锐眼光。你的潜在盲点在于可能过于享受辩论而忽略了维护关系的重要性; 在实际执行和跟进方面存在不足, 容易从兴奋转向无聊; 可能过于关注大局而忽略重要细节; 有时给人不够严肃或不够可靠的印象; 在需要深入专注和持久努力时容易分心。成长方向建议你在辩论中注意对方的感受, 不是每个人都享受智力挑战; 培养执行力和坚持力, 将好的想法变为现实; 学会关注细节, 它们往往是成功的关键; 在追求新奇的同时, 也培养深度和持久性。",
            'ESTJ': "总经理型(ESTJ) - ESTJ是MBTI中最有组织力和领导力的人格类型。你实际、现实主义, 果断, 一旦下决心就会马上行动。你是天生的组织者和管理者, 能够高效地将项目和人组织起来, 用最有效率的方法达成目标。你有强烈的责任感和使命感, 对自己和他人都有很高的标准。你尊重规则和传统, 认为秩序和纪律是高效运转的基础。你直率坦诚, 不喜欢拐弯抹角, 这种品质让你在领导岗位上赢得信任。你有出色的判断力和决策能力, 能够在信息充分的情况下做出快速而准确的选择。你注重结果和效率, 不喜欢浪费时间在无意义的讨论上。你是可靠的执行者, 一旦承诺就会全力以赴。你善于建立和维护系统化的流程, 让一切有条不紊地运行。你有天生的管理才能, 能够清楚地分配任务、设定标准、跟进进度。你的核心优势包括出色的组织和管理能力, 强大的决策力和执行力, 高度的责任感和可靠性, 建立和维护高效系统的能力, 以及在压力下保持冷静和高效的能力。你的潜在盲点在于可能过于注重效率和规则而忽略他人的情感需求; 在做决定时可能过于果断而缺乏充分倾听; 对与自己方式不同的人可能不够包容; 可能过于死板, 难以适应灵活变化的环境; 在表达温柔和关怀方面存在困难。成长方向建议你在做决定之前, 花时间倾听他人的意见和感受; 认识到效率和温情并不矛盾, 温暖的领导方式往往更有效; 学会灵活变通, 并非所有情况都需要按既定规则处理; 有意识地表达对团队成员的认可和感谢。",
            'ESFJ': "执政官型(ESFJ) - ESFJ是MBTI中最温暖、最善于照顾他人的人格类型。你热心肠、有责任心、合作精神强, 希望周围的环境温馨而和谐, 并为此果断地执行。你有天生的社交能力, 能够敏锐地感知到他人的需求, 并主动提供帮助。你重视传统和社群, 努力维护身边人的幸福和团结。你慷慨大方, 乐于为他人付出时间和精力。你有出色的组织能力, 善于安排各种活动和聚会, 让大家都感到愉快。你忠于自己的承诺, 是朋友和家人最可靠的依靠。你注重和谐, 善于调解矛盾, 让不同的人能够在同一团队中融洽相处。你有很强的责任感和使命感, 尤其在照顾家人和帮助朋友方面表现突出。你的温暖和善良是真诚的, 这让身边的人在你面前感到安全和被爱护。你的核心优势包括出色的社交能力和人际协调能力, 强大的同理心和关怀他人的能力, 优秀的组织和活动策划能力, 高度的忠诚度和可靠性, 以及营造和谐氛围的天赋。你的潜在盲点在于可能过于在意他人的看法而忽略自己的真实需求; 在面对冲突时倾向于回避或妥协, 而非正面解决问题; 可能过于依赖他人的认可来获得自我价值感; 对变化和不确定性的适应能力较弱; 有时过于控制, 希望一切都按照自己的预期进行。成长方向建议你学会优先照顾自己的需求, 认识到自我关怀不是自私; 在冲突面前不要退缩, 建设性的冲突能够带来更好的结果; 培养独立思考的能力, 不要总是依赖他人的认可; 接受变化是生活的一部分, 尝试从中找到积极的一面。",
            'ENFJ': "主人公型(ENFJ) - ENFJ是MBTI中最有魅力和影响力的人格类型之一。你热情、为他人着想、易感应、有责任心。你非常注重他人的感情、需求和动机, 天生就懂得如何理解和激励人心。你有卓越的领导才能, 但你的领导方式不是通过权威, 而是通过感召和激励。你善于发现每个人的潜能, 并帮助他们实现最好版本的自己。你有天生的教育和指导能力, 无论是正式的教学还是日常的引导, 你都能让人心服口服。你有强大的同理心, 能够深刻地理解他人的感受和处境。你有明确的价值观和理想, 愿意为了实现这些理想而努力奋斗。你有出色的沟通能力, 无论是公开演讲还是私下交流, 你都能打动人心。你有组织能力, 善于协调不同人的力量来达成共同的目标。你的核心优势包括卓越的领导力和激励他人的能力, 强大的同理心和理解他人情感的能力, 出色的沟通和表达能力, 发现和培养他人潜能的天赋, 以及营造团队凝聚力和归属感的能力。你的潜在盲点在于可能过于关注他人的需求而忽略自己的需要; 对失败和他人的负面反馈过于敏感; 可能过于理想化, 对自己和他人都设定过高的标准; 有时过于热情主动, 给人压迫感; 在需要做不受欢迎的决定时可能犹豫不决。成长方向建议你学会设定边界, 认识到自己的需求同样重要; 接受批评是成长的必经之路, 不要过于在意他人的评价; 适当降低标准, 认识到足够好有时比完美更健康; 给别人留出自主空间, 不是每个人都需要被指导和帮助。",
            'ENTJ': "指挥官型(ENTJ) - ENTJ是MBTI中最有魄力和远见的人格类型。你坦诚、果断, 有天生的领导能力。你能很快看到组织程序和政策中的不合理性和低效能性, 并发展实施有效的全面系统来解决问题。你是天生的战略家, 能够看到别人看不到的大局和可能性。你有强大的意志力和执行力, 一旦确定目标就会全力以赴, 不达目的誓不罢休。你自信、果断, 在别人犹豫不决的时候你已经开始了行动。你有出色的组织和管理能力, 能够将复杂的项目分解为可执行的步骤, 并高效地推进。你重视效率和成果, 不喜欢浪费时间和资源在无意义的事情上。你有强大的沟通能力, 能够在公开场合清晰有力地表达自己的观点。你不畏惧挑战和困难, 反而将它们视为成长的机遇。你的自信心和决断力让你在领导岗位上如鱼得水。你的核心优势包括卓越的战略规划和决策能力, 强大的领导力和影响力, 出色的组织和执行能力, 快速识别和解决问题的能力, 以及在高压力环境下保持高效运转的能力。你的潜在盲点在于可能过于强势和自信, 忽略他人的感受和建议; 在追求效率的过程中可能显得冷酷无情; 可能对自己和他人都设定过高的标准; 在处理复杂情感问题时缺乏耐心和细腻; 可能过于关注结果而忽略了过程中的体验。成长方向建议你学会倾听和尊重不同的意见, 优秀的领导者也需要善于听取声音; 在追求效率的同时关注团队成员的感受和需求; 认识到领导力不仅仅是掌控, 还包括服务和支持; 培养耐心, 学会享受过程而不仅仅是结果。",
        }

        dimensions = []

        e_score = scores.get('E', 0)
        i_score = scores.get('I', 0)
        if e_score > i_score:
            dimensions.append({'dimension': 'E vs I', 'result': 'E', 'score': f'E:{e_score} vs I:{i_score}', 'description': self.dimension_descriptions.get('E', '')})
        else:
            dimensions.append({'dimension': 'E vs I', 'result': 'I', 'score': f'E:{e_score} vs I:{i_score}', 'description': self.dimension_descriptions.get('I', '')})

        s_score = scores.get('S', 0)
        n_score = scores.get('N', 0)
        if s_score > n_score:
            dimensions.append({'dimension': 'S vs N', 'result': 'S', 'score': f'S:{s_score} vs N:{n_score}', 'description': self.dimension_descriptions.get('S', '')})
        else:
            dimensions.append({'dimension': 'S vs N', 'result': 'N', 'score': f'S:{s_score} vs N:{n_score}', 'description': self.dimension_descriptions.get('N', '')})

        t_score = scores.get('T', 0)
        f_score = scores.get('F', 0)
        if t_score > f_score:
            dimensions.append({'dimension': 'T vs F', 'result': 'T', 'score': f'T:{t_score} vs F:{f_score}', 'description': self.dimension_descriptions.get('T', '')})
        else:
            dimensions.append({'dimension': 'T vs F', 'result': 'F', 'score': f'T:{t_score} vs F:{f_score}', 'description': self.dimension_descriptions.get('F', '')})

        j_score = scores.get('J', 0)
        p_score = scores.get('P', 0)
        if j_score > p_score:
            dimensions.append({'dimension': 'J vs P', 'result': 'J', 'score': f'J:{j_score} vs P:{p_score}', 'description': self.dimension_descriptions.get('J', '')})
        else:
            dimensions.append({'dimension': 'J vs P', 'result': 'P', 'score': f'J:{j_score} vs P:{p_score}', 'description': self.dimension_descriptions.get('P', '')})

        return {
            'type': mbti_type,
            'description': type_descriptions.get(mbti_type, 'unknown'),
            'dimensions': dimensions,
            'career_suggestions': self._get_career_suggestions(mbti_type),
            'relationship_advice': self._get_relationship_advice(mbti_type)
        }

    def _get_career_suggestions(self, mbti_type):
        suggestions = {
            'ISTJ': "作为ISTJ型人格, 你天生适合需要精确性、可靠性和系统性思维的职业。你最适合的职业领域包括财务管理、审计与合规、行政管理、工程技术和法律等需要严谨态度的行业。具体推荐的岗位有注册会计师、审计师、税务专员、项目经理、系统管理员、数据库管理员、质量检测工程师、法律顾问、法官助理、行政管理主管。在职业发展中, 你的优势在于能够建立和维护高效的流程系统, 确保工作质量和一致性。建议你发挥组织能力强的特点, 逐步向管理层发展, 成为团队中稳定可靠的核心人物。需要注意的是, 在职场中你可能对变化和新技术的接受速度较慢, 建议有意识地保持学习的开放态度。你最适合的工作环境是有明确规则和流程、重视质量和效率、能够独立完成深度工作。避免过于模糊或需要频繁社交应酬的岗位。你的职业成长路径通常是从技术专家到管理岗位的渐进式发展。在职业规划方面, 建议你制定清晰的3-5年职业发展路线图, 明确每个阶段的目标和需要的技能提升。可以考虑获取相关专业认证来增强竞争力, 比如CPA、PMP等。在职场沟通中, 虽然你不是天生的演说家, 但通过准备充分的结构化表达, 你同样能有效传达观点。建议参加一些演讲和沟通技巧的培训课程。在团队管理方面, 你可以尝试带新人或实习生, 这将锻炼你的指导能力和耐心。同时, 保持对新技术的敏感度, 定期参加行业研讨会, 了解最新趋势。你的职业发展潜力很大, 关键在于持续学习和适度开放。",
            'ISFJ': "作为ISFJ型人格, 你最适合能够帮助他人、营造和谐环境的职业。你天生具有关怀他人的天赋, 在需要耐心、细心和同理心的领域会大放异彩。推荐的职业领域包括医疗护理、教育培训、社会服务、人力资源和客户关系管理。具体推荐的岗位有注册护士、特殊教育教师、学校心理咨询师、社工、人力资源专员、客户关系经理、行政助理、图书馆管理员、室内设计师、营养师。在职业发展中, 你的优势在于出色的观察力和对细节的关注, 以及对他人需求的敏锐感知。建议你寻找能够直接帮助他人的岗位, 这将带给你最大的职业满足感。需要注意的是, 你可能过于牺牲自己来满足工作需求, 建议学会设定合理的工作边界。你最适合的工作环境是温暖友善、团队氛围和谐、工作节奏稳定有序。你的职业成长路径可以通过深耕专业领域, 成为受人尊敬的行业专家。在职业规划方面, 建议你注重专业能力的深度发展, 在一个领域成为真正的专家。可以考虑获取心理咨询、教育、医疗等相关领域的专业资质。在职场沟通中, 你的细心和体贴是独特的优势, 善于关注团队成员的情绪变化。建议学习一些基本的冲突管理和谈判技巧。在团队协作方面, 你可以担任协调者和沟通桥梁的角色, 发挥你理解不同人需求的能力。同时, 不要忘记为自己争取应得的认可和晋升机会, 你的贡献值得被看见。在工作与生活平衡方面, 设定明确的界限, 避免过度付出导致职业倦怠。你的职业道路是稳步上升的, 坚持专业深耕是关键。",
            'INFJ': "作为INFJ型人格, 你最适合能够发挥洞察力、创造力和帮助他人成长的职业。你天生理解人性, 善于发现深层次的问题和意义。推荐的职业领域包括心理咨询、教育培训、写作与出版、艺术创作、人力资源和组织发展。具体推荐的岗位有心理咨询师、心理治疗师、职业规划师、大学教授(人文社科领域)、作家、编辑、艺术治疗师、人力资源发展专员、非营利组织管理者、品牌策划师。在职业发展中, 你的优势在于深刻的洞察力、出色的写作和沟通能力, 以及激励他人的天赋。建议你选择有使命感的工作, 这将激发你最大的潜能。需要注意的是, 你可能对工作环境的精神氛围要求很高, 建议在求职时充分了解企业文化。你最适合的工作环境是有意义、有价值、能够发挥创造力的地方。避免过于机械化或缺乏人文关怀的岗位。在职业规划方面, 建议你寻找与你核心价值观高度一致的岗位和组织文化。企业文化对INFJ来说比薪资更重要。可以考虑结合写作和咨询的复合型职业路径。在职场沟通中, 你的深度思考能力让你擅长撰写高质量的分析报告和提案。建议学习一些项目管理和执行方法论, 让你的远见能够落地。在职业发展中, 不要害怕跳槽寻找真正契合的环境。你的直觉往往能帮助你做出正确的职业选择。建议保持至少一项创造性爱好, 这将帮助你在高压工作中保持心理平衡。你的职业潜力在于能够将深刻的洞察转化为有价值的行动方案。",
            'INTJ': "作为INTJ型人格, 你最适合需要战略思维、独立判断和创新能力的职业。你天生善于设计和优化系统, 在复杂的智力挑战中能够找到最优解。推荐的职业领域包括科学研究、技术研发、战略规划、金融分析和高级咨询。具体推荐的岗位有科学家、研究员、系统架构师、首席技术官、战略咨询顾问、投资分析师、大学教授(理工科领域)、专利工程师、数据科学家、企业家。在职业发展中, 你的优势在于卓越的规划能力、快速学习和掌握复杂系统的能力, 以及独立高效的执行力。建议你寻找能够自主决策、有挑战性的岗位, 这将充分发挥你的潜能。需要注意的是, 在团队协作中要注意沟通方式, 多倾听同事的想法。你最适合的工作环境是能够独立思考、有充分自主权、以能力和成果为导向。避免过于官僚化或限制创新的岗位。你的职业发展空间广阔, 适合向高级管理和专业领袖方向发展。在职业规划方面, 建议你明确自己的长期职业愿景, 并制定分阶段的实施计划。你的战略思维能力让你非常适合高级管理或专业领袖岗位。可以考虑MBA或其他高级管理课程来补充商业知识。在职场沟通中, 建议练习用更通俗的语言向非技术背景的人解释复杂概念, 这将大大提升你的影响力。在团队领导方面, 学会根据团队成员的不同特点分配任务和激励方式。同时, 建议培养至少一个非工作相关的兴趣, 这将帮助你保持思维的新鲜度。你的职业天花板很高, 关键在于持续提升人际影响力。",
            'ISTP': "作为ISTP型人格, 你最适合需要动手能力、技术理解和灵活应变的职业。你天生擅长理解和操作物理系统, 在需要实际操作的环境中如鱼得水。推荐的职业领域包括工程技术、制造业、执法与安全、体育运动和医疗技术。具体推荐的岗位有机械工程师、电气工程师、飞行员、消防员、刑侦分析师、外科医生、汽车技师、软件开发工程师(尤其是系统级)、数据分析员、运动教练。在职业发展中, 你的优势在于出色的动手能力、冷静的危机处理能力和快速学习能力。建议你选择能够不断接触新技术、解决实际问题的岗位。需要注意的是, 你可能对常规性和重复性工作缺乏耐心, 建议寻找有变化和挑战的工作内容。你最适合的工作环境是灵活自主、以结果为导向、能够实际动手操作。避免需要大量文书工作或复杂社交应酬的岗位。在职业规划方面, 建议你选择一个允许你不断学习新技能和解决新问题的领域。技术类工作最适合你, 但也要注意培养沟通和团队协作能力。可以考虑考取相关专业执照或认证来增强职业竞争力。在职场沟通中, 建议你在做出技术判断后, 多花时间向团队成员解释你的 reasoning 过程。在职业发展中, 不要局限于一个狭窄的技术领域, 跨领域的经验将让你更有价值。建议定期更新你的技能组合, 跟上行业技术发展的步伐。你的适应能力和动手能力是核心竞争力, 找到能够充分发挥这些特长的平台将让你如虎添翼。",
            'ISFP': "作为ISFP型人格, 你最适合能够发挥审美天赋、创意才能和关怀之心的职业。你对美和品质有着天生的敏感度, 在创意和助人领域都能找到满足感。推荐的职业领域包括艺术设计、时尚美容、医疗健康、教育培训和社会服务。具体推荐的岗位有平面设计师、室内设计师、时装设计师、摄影师、插画师、美容师、按摩治疗师、营养师、幼儿教育教师、宠物护理师、园艺师、花艺设计师。在职业发展中, 你的优势在于出色的审美能力、细致的手工技能和真诚的人际关系。建议你选择能够发挥创造力和个性表达的工作, 这将让你充满激情。需要注意的是, 你可能需要额外努力来应对商业和管理方面的工作, 建议学习基本的项目管理技能。你最适合的工作环境是灵活自由、有创意空间、环境优美舒适。避免过于竞争激烈或需要不断社交应酬的岗位。在职业规划方面, 建议你建立一个作品集或项目集来展示你的创意能力, 这比任何简历都更有说服力。可以考虑参加设计比赛、艺术展览等活动来提升曝光度。在职场沟通中, 学会用视觉化的方式表达你的创意理念, 比如制作设计稿或样品。在职业发展中, 不要害怕尝试不同领域的创意工作, 你的审美天赋在多个行业都有价值。建议学习一些基础的商业和营销知识, 这将帮助你更好地推广自己的作品和服务。同时, 寻找与你风格互补的合作伙伴, 互补的团队能够创造出更优秀的成果。你的创造力是稀缺资源, 善用这份天赋。",
            'INFP': "作为INFP型人格, 你最适合能够表达创造力、实现理想和帮助他人的职业。你的理想主义和创造力是你最大的职业资产。推荐的职业领域包括文学创作、心理咨询、艺术设计、教育公益和社会创新。具体推荐的岗位有作家、诗人、编剧、心理咨询师、艺术治疗师、大学教授(人文学科)、社会工作者、非营利组织项目专员、翻译、编辑、UX设计师、品牌文案策划。在职业发展中, 你的优势在于深厚的文字功底、丰富的想象力和对人类情感的深刻理解。建议你寻找与你价值观一致的工作, 这会给你带来持久的动力。需要注意的是, 你可能需要加强执行力, 将创意转化为具体的成果。建议学习项目管理和时间管理的基本方法。你最适合的工作环境是自由宽松、有创意空间、价值观一致。避免过于商业化或与你道德理念冲突的岗位。在职业规划方面, 建议你寻找能够将你的价值观、创造力和专业技能结合起来的工作。纯粹商业化或与你信念冲突的工作会让你感到空虚。可以考虑自由职业或创业, 这将给你更大的创作自由。在职场沟通中, 你的文字表达能力是你的强项, 善用写作来传达你的想法和立场。在职业发展中, 学习基本的商业运营和财务管理知识, 这将帮助你的创意项目可持续发展。建议建立一个创作习惯, 每天固定时间进行创作或写作。同时, 寻找一个支持你创作的社区或圈子, 与志同道合的人交流将激发你的创作灵感。你的理想主义是最好的驱动力, 找到合适的表达渠道。",
            'INTP': "作为INTP型人格, 你最适合需要深度分析、逻辑推理和创新思维的职业。你的大脑天生适合处理复杂的理论和系统问题。推荐的职业领域包括科学研究、技术开发、数据分析、学术研究和哲学思辨。具体推荐的岗位有软件架构师、数据科学家、AI研究员、大学教授(理工科和哲学)、理论物理学家、数学家、经济学家、专利分析师、技术顾问、网络安全专家。在职业发展中, 你的优势在于卓越的逻辑分析能力、快速学习和抽象思维能力。建议你选择能够持续学习和挑战智力的岗位, 避免重复性的工作。需要注意的是, 你可能需要加强沟通和团队协作能力, 这将在职业发展中发挥越来越重要的作用。你最适合的工作环境是自主性强、智力挑战高、尊重独立思考。避免需要大量社交互动或执行常规性任务的岗位。你的职业发展前景广阔, 在技术前沿领域尤其有优势。在职业规划方面, 建议你选择处于技术前沿的领域, 你的学习和理解能力在快速发展的行业中具有巨大优势。可以考虑参与开源项目或技术社区来积累声誉和经验。在职场沟通中, 建议你在提出创新想法的同时, 也提供一个初步的实施方案, 这将大大提高你的说服力。在职业发展中, 不要满足于理论研究, 也要关注技术的实际应用和商业化可能性。建议定期参加技术大会或学术研讨会, 保持对行业趋势的了解。同时, 培养一些软技能, 如演讲、写作和团队管理, 这些技能将在你向更高层次发展时变得至关重要。你的创新能力是核心竞争力, 找到合适的平台充分发挥。",
            'ESTP': "作为ESTP型人格, 你最适合需要快速反应、实际行动和人际交往能力的职业。你的冒险精神和应变能力是独特的职业优势。推荐的职业领域包括商业销售、企业管理、体育竞技、执法安全和创业创新。具体推荐的岗位有销售总监、企业家、房地产经纪人、营销经理、运动员、体育教练、消防员、刑警、急诊医生、飞行员、活动策划师、股票交易员。在职业发展中, 你的优势在于出色的应变能力、天生的社交魅力和强大的执行力。建议你选择节奏快、充满挑战的岗位, 你会在高压环境下表现出色。需要注意的是, 你可能需要培养长期规划和耐心的品质, 这有助于你在职业生涯中走得更远。你最适合的工作环境是充满活力、变化多样、以成果为导向。避免需要长时间独处或进行深度理论研究的岗位。你的职业发展空间很大, 适合向创业和高级管理方向发展。在职业规划方面, 建议你选择节奏快、变化多、奖励机制明确的工作环境。你的行动力和应变能力在创业和销售领域有天然优势。可以考虑学习一些商业管理知识来补充你的实战经验。在职场沟通中, 你的直率和魅力是你的优势, 但也要学会根据不同对象调整沟通方式。在职业发展中, 建议你制定中期目标, 比如三年内达到某个职位或收入水平, 这将帮助你保持方向感。同时, 注意培养耐心和坚持力, 一些最有价值的成就需要长期积累。建议在精力充沛的年轻阶段多尝试不同领域, 找到最适合自己的发展方向。你的活力和执行力是巨大优势, 找到合适的赛道将让你脱颖而出。",
            'ESFP': "作为ESFP型人格, 你最适合能够发挥社交能力、创造力和表现力的职业。你的热情和感染力是你最大的职业资产。推荐的职业领域包括娱乐传媒、时尚美业、旅游酒店、教育培训和销售服务。具体推荐的岗位有演员、主持人、活动策划师、公关经理、旅游顾问、酒店管理、销售经理、时尚买手、健身教练、幼教老师、婚礼策划师、美食博主。在职业发展中, 你的优势在于出色的社交能力、即兴表现能力和让周围人感到愉快的天赋。建议你选择能够与人互动、展示创意的岗位。需要注意的是, 你可能需要加强时间管理和组织规划能力。你最适合的工作环境是充满社交互动、氛围轻松愉快、能够发挥个人魅力。避免需要长时间独处或高度结构化的岗位。你的职业发展路径可以沿着社交影响力和创意表现力的方向不断提升。在职业规划方面, 建议你选择能够与人频繁互动、发挥个人魅力的岗位。你的社交能力和感染力在服务、教育和娱乐行业有巨大价值。可以考虑学习一些市场营销和品牌管理的知识来拓展职业空间。在职场沟通中, 你的即兴表达能力很强, 但在正式场合建议提前准备要点, 这样既能保持自然又有条理。在职业发展中, 建议你积累一些可以量化的工作成果, 比如客户满意度数据或活动参与人数, 这将在你寻求晋升或加薪时很有帮助。同时, 不要忽视持续学习的重要性, 在享受工作的同时也要提升专业技能。你的热情和亲和力是核心竞争力, 找到能够充分发挥这些特质的平台。",
            'ENFP': "作为ENFP型人格, 你最适合能够发挥创造力、沟通能力和激励他人天赋的职业。你的热情和创新精神在任何领域都是宝贵的资产。推荐的职业领域包括创意产业、教育培训、媒体传播、咨询策划和社会创新。具体推荐的岗位有创意总监、品牌策划师、记者、主持人、创业顾问、培训师、大学教授、公关经理、心理辅导员、编剧、社交媒体运营总监、产品经理。在职业发展中, 你的优势在于卓越的创意能力、强大的沟通感染力和发现新机会的眼光。建议你选择能够不断接触新事物、发挥创造力的岗位。需要注意的是, 你可能需要培养专注力和执行力, 学会将好的想法落地实现。你最适合的工作环境是充满创意自由、鼓励创新、团队氛围积极向上。避免过于官僚化或限制创意的岗位。你的职业发展路径可以从创意岗位逐步发展为具有影响力的领导者。在职业规划方面, 建议你选择允许创新和变革的岗位和组织。你的创意能力和激励他人的天赋在创业公司、咨询公司和创意产业有巨大价值。可以考虑建立个人品牌来扩大影响力。在职场沟通中, 你的感染力很强, 但要学会根据受众调整信息的深度和形式, 对管理层要用数据说话, 对客户要用故事打动人心。在职业发展中, 建议你至少在一个领域深耕三年以上, 建立扎实的专业基础后再考虑拓展。同时, 培养一些项目管理的基本功, 比如制定里程碑、追踪进度等, 这将帮助你的创意想法更好地落地。你的创意和热情是稀缺资源, 找到合适的平台和团队来最大化发挥。",
            'ENTP': "作为ENTP型人格, 你最适合需要创新思维、快速学习和解决复杂问题的职业。你的辩论精神和系统分析能力是独特的职业优势。推荐的职业领域包括创业创新、技术咨询、法律辩护、媒体评论和战略规划。具体推荐的岗位有企业家、创业顾问、管理咨询师、律师、记者、评论员、产品经理、风险投资人、市场营销策略师、政治分析师、大学教授、发明家。在职业发展中, 你的优势在于卓越的系统分析能力、创新思维和出色的辩论沟通能力。建议你选择充满智力挑战和变化的工作, 你会在解决难题中找到最大的满足。需要注意的是, 你可能需要加强执行力和持续跟进的能力, 避免三分钟热度。你最适合的工作环境是充满智力挑战、鼓励创新、尊重独立思考。避免过于传统或限制创新空间的岗位。你的职业发展前景广阔, 创业和技术创新领域尤其适合你。在职业规划方面, 建议你选择充满挑战和变化的工作环境, 你的创新能力和分析思维在创业和咨询领域有巨大优势。可以考虑建立一个多元化的职业网络, 这将为你带来更多的机会和资源。在职场沟通中, 你的辩论能力很强, 但要学会区分什么时候该辩论、什么时候该倾听, 这将大大提升你的领导力。在职业发展中, 建议你选择一个能够看到从创意到落地全过程的岗位, 完整的参与感对你很重要。同时, 培养执行力和坚持力, 你的很多好想法都需要更长的实施周期才能见效。建议找到一位能够互补的合作伙伴, 对方能帮你把控细节和跟进进度。你的创新精神是核心竞争力, 找到合适的舞台尽情发挥。",
            'ESTJ': "作为ESTJ型人格, 你最适合需要组织管理、执行力和领导力的职业。你的决断力和组织能力是天然的领导优势。推荐的职业领域包括企业管理、行政管理、法律司法、财务管理和教育培训。具体推荐的岗位有项目经理、运营总监、首席执行官、法官、律师、财务总监、学校校长、军队军官、供应链管理经理、行政主管、质量管理经理。在职业发展中, 你的优势在于出色的组织和管理能力、强大的决策力和执行力。建议你发挥领导才能, 逐步向更高级的管理岗位发展。需要注意的是, 在管理中要关注团队成员的情感需求, 学会在效率和人文关怀之间取得平衡。你最适合的工作环境是有明确目标和结构、重视效率和成果、以团队协作为核心。避免过于模糊或缺乏组织结构的工作环境。你的职业发展路径通常是从业务骨干到团队管理再到高级领导的阶梯式上升。在职业规划方面, 建议你制定清晰的3-5年职业发展路线图, 从业务骨干到中层管理再到高级领导层逐步晋升。你的组织能力和执行力让你天然适合管理岗位。可以考虑获取MBA或其他管理类资质来增强理论功底。在职场沟通中, 学会在保持效率的同时增加一些温度, 用关心和认可来激励团队成员。在职业发展中, 建议你跨部门轮岗以获得更全面的业务理解。同时, 关注行业趋势和技术变革, 不要让对传统方式的坚持让你落后于时代。建议培养一个你感兴趣的副业或爱好, 这将帮助你在高强度工作之余保持心理平衡。你的领导力是核心竞争力, 持续提升将让你走得更远。",
            'ESFJ': "作为ESFJ型人格, 你最适合能够发挥社交能力、关怀精神和组织能力的职业。你天生擅长建立和维护人际关系, 在需要团队合作和客户服务的领域表现突出。推荐的职业领域包括医疗护理、教育培训、人力资源、客户服务和公共关系。具体推荐的岗位有护士、学校教师、幼儿园园长、人力资源经理、客户服务主管、公关经理、社会工作主管、酒店管理、活动策划师、销售经理、社区服务主任、健康顾问。在职业发展中, 你的优势在于出色的人际协调能力、强大的同理心和优秀的组织能力。建议你选择能够直接与人互动、帮助他人的岗位。需要注意的是, 你可能需要学会在帮助他人的同时维护自己的权益和需求。你最适合的工作环境是团队合作密切、氛围和谐温馨、重视人际关系。避免需要长时间独处或竞争过于激烈的岗位。你的职业发展可以从专业岗位向管理和服务领导方向发展。在职业规划方面, 建议你选择以人为本的工作环境, 你的关怀能力和组织天赋在人力资源、教育和客户服务领域有巨大价值。可以考虑考取人力资源管理师或心理咨询师等认证来提升专业度。在职场沟通中, 你的倾听能力很强, 但也要学会清晰表达自己的立场和需求。在职业发展中, 建议你主动争取项目管理或团队领导的机会, 你其实拥有出色的领导潜质。同时, 学习一些数据分析和决策工具, 这将帮助你在关心人的同时也能做出客观的判断。建议定期反思自己的职业满意度, 不要因为害怕冲突就一直留在舒适区。你的温暖和组织能力是核心竞争力, 大胆展现你的领导潜能。",
            'ENFJ': "作为ENFJ型人格, 你最适合能够发挥领导力、激励他人和沟通能力的职业。你天生能够感召和鼓舞人心, 在教育、咨询和领导领域有巨大优势。推荐的职业领域包括教育培训、心理咨询、人力资源管理、非营利管理和公共演讲。具体推荐的岗位有大学教授、企业培训师、心理咨询师、人力资源总监、学校校长、非营利组织执行官、公关总监、职业教练、政治家、社区领袖、组织发展顾问。在职业发展中, 你的优势在于卓越的领导魅力、深刻的同理心和出色的沟通能力。建议你选择能够影响和帮助他人成长的岗位, 这将带给你最大的职业满足感。需要注意的是, 你可能需要学会设定个人边界, 避免过度投入他人事务而忽略自己。你最适合的工作环境是重视人文关怀、鼓励团队合作、有积极向上的文化。避免冷漠或竞争过度的工作环境。你的职业发展路径可以从教育和咨询领域向更高层次的领导和影响力角色发展。在职业规划方面, 建议你选择能够发挥激励和培养他人能力的岗位, 你的领导魅力和同理心在教育、培训和人力资源高管岗位有巨大价值。可以考虑建立导师关系网络来扩展你的影响力。在职场沟通中, 你的感染力很强, 但也要学会设定边界, 不要让帮助他人的冲动消耗了你的精力。在职业发展中, 建议你逐步从执行者角色转向战略领导者角色, 你的大局观和人文关怀将在更高层次上发挥价值。同时, 学习一些财务和运营管理的基本知识, 这将让你成为更全面的领导者。建议在领导风格上不断精进, 在感召力和决断力之间找到平衡。你的影响力是核心竞争力, 找到更大的舞台来施展你的领导才华。",
            'ENTJ': "作为ENTJ型人格, 你最适合需要战略眼光、领导力和执行力的职业。你的决断力和远见卓识使你成为天生的领导者。推荐的职业领域包括高级管理、创业创新、战略咨询、法律金融和政治领域。具体推荐的岗位有首席执行官、首席运营官、战略咨询总监、管理咨询合伙人、投资银行家、企业律师、政治家、军事指挥官、创业者、风险投资人、商学院教授、高管教练。在职业发展中, 你的优势在于卓越的战略规划能力、强大的领导力和影响力、以及高效的执行力。建议你选择有挑战性、能够发挥领导才能的岗位, 你会在管理和发展团队中找到最大的满足。需要注意的是, 在领导过程中要关注团队成员的感受, 培养更柔软的领导方式。你最适合的工作环境是目标明确、竞争激烈、以成果为导向、能够自主决策。避免缺乏挑战或限制领导空间的工作。你的职业发展空间广阔, 适合向最高管理层和创业方向发展。在职业规划方面, 建议你瞄准最高管理层或创业方向, 你的战略思维和执行力让你天生就是做大事的人。可以考虑参与董事会或行业组织来扩大你的影响力。在职场沟通中, 你的说服力很强, 但要记住真正的领导力不仅在于让人听从, 更在于让人心悦诚服。在职业发展中, 建议你培养继任者, 一个优秀的领导者应该能够培养出更多的领导者。同时, 不要忽视对新兴技术和商业模式的关注, 保持学习和适应的能力将让你在快速变化的环境中保持领先。建议在追求业绩的同时也注重企业文化建设, 一个有凝聚力的团队将让你的战略执行更加高效。你的领导力是核心竞争力, 持续挑战更高的目标。",
        }
        return suggestions.get(mbti_type, '')

    def _get_relationship_advice(self, mbti_type):
        advice = {
            'ISTJ': "在人际交往中, 你是一个可靠、稳定但内敛的人。你用实际行动来表达关心和爱意, 而不是甜言蜜语。你可能不会经常说爱你, 但你会记住伴侣的每个重要日子, 默默做好一切事情。你的忠诚度极高, 一旦认定一个人就会全心全意对待。在友谊方面, 你可能朋友不多, 但每一段友谊都非常深厚和持久。与不同类型的人相处时, 建议你尝试理解情感型人士(F型)的表达方式, 他们可能更需要言语上的确认和情感上的回应。面对直觉型(N型)朋友时, 试着欣赏他们的想象力和抽象思维, 即使你觉得有些不切实际。与同样务实的S型伙伴在一起时, 你们会非常合拍, 能够建立起深厚而稳固的关系。在亲密关系中, 你倾向于通过提供稳定的生活、解决实际问题来表达爱。建议你学会更多地用语言和肢体表达感受, 这对你的伴侣来说非常重要。尝试在适当的时候分享你的内心想法, 不要总是把一切都藏在心里。你的伴侣需要感受到你不仅仅是可靠的, 更是温暖的。在处理矛盾时, 学会先关注对方的情绪, 再讨论解决问题的方案。定期安排一些浪漫的约会和惊喜, 这将让你的关系更加丰富和美好。在日常相处中, 建议你尝试一些新的互动方式, 比如一起尝试新的活动或去新的地方旅行, 这将给关系带来新鲜感。学会在关系中也放松规则, 偶尔的即兴和随性会让你的伴侣感到惊喜。当伴侣有负面情绪时, 不要急于给出解决方案, 有时候他们只是需要一个倾听者。记住, 在亲密关系中, 展现脆弱不是软弱, 而是信任的表现。你可以试着与伴侣分享你工作中的挑战和困扰, 这将让你们的关系更加深入和真实。你的忠诚和可靠是关系的基石, 在此基础上增加温暖和柔软, 你将成为最理想的伴侣。",
            'ISFJ': "在人际交往中, 你是最温暖、最体贴的人。你总是第一个注意到朋友不开心的人, 也是最愿意为他人付出的人。你的人际关系以深度和质量为特点, 你重视每一个在你生命中的人。你对朋友和家人非常忠诚, 会记住他们说的每一句话、每一个偏好。你的善良和体贴让你成为朋友圈中的暖心人, 大家遇到困难时第一个想到的就是你。与不同类型的人相处时, 你需要意识到不是所有人都像你一样善于体察他人。思考型(T型)的人可能需要更直接的沟通方式, 他们不一定能读懂你的暗示。在面对外向型(E型)朋友时, 要记得照顾自己的社交电量, 适度表达需要独处的需求。与同样感性的人在一起时, 你们会建立非常深厚的情感连接。在亲密关系中, 你通过无微不至的照顾来表达爱。建议你也要学会接受他人的帮助, 不要总是把自己放在最后一位。在关系中表达自己的需求和不满同样重要, 压抑情绪只会让问题积累。当感到被忽视或不被感激时, 要勇敢地与伴侣沟通。你的付出值得被看见和感恩。建议定期与伴侣进行深入的交流, 分享你的内心感受和需求。在日常相处中, 建议你培养表达正面感受的习惯, 不只是在对方需要帮助时才出现, 也要在对方做得好的时候给予真诚的赞美。学会接受伴侣的照顾和帮助, 这不是软弱的表现, 而是让关系更加平衡健康的方式。当遇到矛盾时, 不要总是退让和忍让, 你真实的想法和需求同样值得被重视。建议你与伴侣一起建立一些共同的兴趣和活动, 这将增加你们之间的连接和默契。记住, 一个健康的关系需要双方都感到满足和被重视, 不要总是把自己放在最后。你的温暖和体贴是珍贵的, 让它成为关系的纽带而非你的枷锁。",
            'INFJ': "在人际交往中, 你追求的是深层次的灵魂连接。你可能朋友不多, 但每一段关系都充满了深度和意义。你有着惊人的直觉, 常常能读懂别人自己都没意识到的情感。你的同理心让你成为优秀的倾听者和建议者。但同时, 你对关系的要求也很高 - 你渴望的是真正的理解和共鸣, 而非表面的社交。与不同类型的人相处时, 你需要理解不是每个人都追求深度的关系。感觉型(S型)的人可能更注重实际和当下, 建议你尝试欣赏他们脚踏实地的品质。面对思考型(T型)的人时, 理解他们的理性表达方式, 他们关心你只是方式不同。与同样直觉型的人在一起时, 你们会产生奇妙的化学反应, 但也要注意不要一起陷入过度理想化。在亲密关系中, 你需要大量的独处时间来恢复能量, 这是完全正常的。建议你提前与伴侣沟通这一需求, 避免让他们误解为疏远。学会在理想与现实之间找到平衡, 不要因为伴侣没有达到你心中的完美标准就感到失望。你的深度和洞察力是珍贵的礼物, 用它来理解和滋养关系, 但也要记得享受关系中简单而平凡的快乐。在日常相处中, 建议你学会享受简单平凡的相处时光, 不是每次交流都需要有深度和意义。一起散步、做饭、看电影, 这些看似普通的时刻同样是关系的重要组成部分。当你需要独处时间时, 用温和的方式表达, 比如我需要一些安静的时间来充电, 一小时后我会来找你。这比突然消失或沉默更容易被理解。建议你对伴侣保持好奇心, 即使已经很了解了, 也要不断发现对方身上新的闪光点。在关系中, 你的深度和洞察力是珍贵的, 但不要用它来过度分析对方的行为和动机。学会信任和放松, 让关系自然而然地发展。你的理想主义是美好的, 但也要学会欣赏当下的幸福。",
            'INTJ': "在人际交往中, 你是一个理性、独立但深度思考的人。你对人际关系有着独特的高标准 - 你更看重智力上的共鸣和相互尊重, 而非表面的友好。你可能不容易与人建立深层连接, 但一旦建立了, 你会是非常忠诚和值得信赖的朋友。你用逻辑而非情感来处理人际问题, 这在某些场合是优势, 但也可能导致你忽略了他人的情感需求。与不同类型的人相处时, 建议你对情感型(F型)的人多一些耐心和同理心, 他们处理问题的方式与你不同但同样有效。面对外向型(E型)朋友时, 理解他们需要社交互动来获取能量, 即使你觉得这没有必要。与同样思考型的人在一起时, 你们能够进行深入的智力交流, 这是你最享受的社交方式。在亲密关系中, 你通过解决问题和提供智慧支持来表达爱。建议你也要学会表达温柔和情感, 这对你的伴侣来说至关重要。在冲突中, 先处理情感再处理问题 - 即使你认为问题本身更重要, 你的伴侣需要先感受到被理解和被关注。你的独立性和深度是你的魅力, 但在关系中也要学会展现脆弱的一面。定期与伴侣分享你的内心世界和未来愿景, 让他们感受到自己是你生命中重要的一部分。在日常相处中, 建议你定期安排一些纯粹放松的约会时间, 不讨论工作、不解决问题, 只是享受彼此的陪伴。这看似浪费效率, 但实际上是为关系的长期健康进行投资。当伴侣向你倾诉烦恼时, 先问一句你需要我帮你分析还是只是想听你说, 这将避免很多不必要的误解。建议你对伴侣的兴趣和爱好表现出真诚的好奇, 即使你不一定感兴趣, 这种尊重和关注对伴侣来说意义重大。在关系中, 你的智慧和远见是宝贵的, 但也要学会活在当下, 享受此刻的温馨和宁静。记住, 最深的爱不是理解对方的一切, 而是在不理解的时候依然选择陪伴和信任。",
            'ISTP': "在人际交往中, 你是一个随性、独立、实际的人。你不喜欢过度复杂的人际关系, 更倾向于少说多做的交往方式。你的幽默感和在危机中冷静沉着的表现让你在很多社交场合中都受到欢迎。你对自己的空间和自由非常重视, 不喜欢被束缚或控制。你用行动而非言语来表达关心, 比如帮朋友修好一辆自行车, 比说一百句关心的话更符合你的风格。与不同类型的人相处时, 理解情感型(F型)的人可能需要更多的言语确认和情感回应。面对判断型(J型)朋友时, 你可能需要更多地遵守约定的时间计划。与同样感知型(P型)的人在一起时, 你们能够轻松自在地相处, 不需要太多的规矩和束缚。在亲密关系中, 你的实际能力和冷静态度是你独特的魅力。建议你学会用言语和肢体语言表达爱意, 不要让行动成为唯一的表达方式。在关系中保持一定的计划性和承诺度, 你的伴侣需要感受到稳定和安全感。当伴侣分享情感时, 认真倾听而不是急于提供解决方案。学会享受安静温馨的相处时光, 不一定要通过冒险和行动来维系感情。定期表达你的感受, 即使这对你说来不太自然。在日常相处中, 建议你主动发起一些不需要冒险的温馨活动, 比如一起做饭、看电影或在家里度过一个安静的周末。这些简单的时刻同样能增进你们的感情。当伴侣表达情感需求时, 不要用行动来代替言语回应, 有时候一句我理解你或者我在乎你比任何实际帮助都更有力量。建议你对伴侣保持一定的时间承诺, 比如每周固定的约会之夜, 这将给伴侣带来安全感和稳定感。在关系中, 你的冷静和能力是迷人的, 但也要学会展现你的柔软面。当你遇到困难或感到脆弱时, 向伴侣倾诉不是示弱, 而是加深你们之间信任和亲密感的方式。你的实际和可靠是关系的支柱, 在此基础上增加情感表达, 你将成为最理想的伴侣。",
            'ISFP': "在人际交往中, 你是一个温和、敏感、真诚的人。你有天生的亲和力, 让人感到放松和舒适。你不善言辞, 但你的善良和体贴通过细微的举动传达给身边的人。你重视真实和深度的关系, 不喜欢虚伪和表面化的社交。你有很强的审美感受力, 经常能通过小细节让身边的人感到被关怀 - 一束精心挑选的花、一张手写的卡片、一顿用心准备的晚餐。与不同类型的人相处时, 理解思考型(T型)的人可能更直接, 他们不是在故意伤害你。面对外向型(E型)的朋友时, 学会在需要独处时勇敢表达, 这不是拒绝而是自我保护。与同样感性的人在一起时, 你们能够建立深厚的情感纽带, 相互理解和支持。在亲密关系中, 你通过温柔的陪伴和贴心的举动来表达爱。建议你也要学会直接表达自己的需求和感受, 不要总是默默期待对方能读懂你的心思。在面对冲突时, 不要总是选择回避, 适度的正面沟通会让关系更加健康。你的伴侣需要了解你的真实想法和感受。学会说不, 在关系中也维护自己的边界。你的温柔和真诚是珍贵的, 用它来建立一段既浪漫又真实的关系。在日常相处中, 建议你主动分享你的内心感受和想法, 不要总是等对方来询问或猜测。你丰富的内心世界是你最大的魅力之一, 让伴侣有机会了解和欣赏它。当你们之间出现分歧时, 不要为了避免冲突而默默退让, 温和而坚定地表达你的立场会让关系更加健康。建议你与伴侣一起创造一些属于你们的艺术体验, 比如一起画画、摄影或听音乐会, 这将加深你们的情感连接。在关系中, 你的温柔和真诚是无价的, 但也要学会接受伴侣对你的关心和帮助。你的敏感不是缺点, 而是一种深刻的感知能力, 用它来感受和回应伴侣的爱意, 同时也用它来保护自己的感受和需求。",
            'INFP': "在人际交往中, 你是一个真诚、温暖、有深度的人。你对人有着天然的信任和善意, 总是愿意看到每个人最好的一面。你追求的是灵魂层面的连接 - 你想要了解一个人的故事、梦想和内心世界。你可能是朋友圈中的心灵导师, 大家都会来找你倾诉和寻求建议。你有出色的倾听能力, 能够给予他人真正被理解和接纳的感受。与不同类型的人相处时, 理解实感型(S型)的人可能更关注实际细节而非抽象意义, 这并不意味着他们缺乏深度。面对思考型(T型)的人时, 接受他们的理性反馈, 这不是对你感受的否定。与同样直觉型的人在一起时, 你们能够进行深入的灵魂对话, 但要注意不要一起过度理想化而脱离现实。在亲密关系中, 你的理想主义可能让你在初期对伴侣抱有过高的期望。建议你接受真实的人都是有缺点的, 爱一个人不仅是欣赏优点, 更是接纳不完美。学会在失望时与伴侣沟通, 而不是默默退缩。你的深情和忠诚是宝贵的品质, 但也要记得在关系中保持自己的独立性。定期与伴侣分享你的内心世界和价值观, 共同构建有意义的生活。你的理想主义在关系中是一束光, 用它照亮你们共同前行的道路。在日常相处中, 建议你与伴侣分享你的梦想和理想, 而不只是你的担忧和不满。你的理想主义是美丽的, 让伴侣有机会了解你内心最珍贵的部分。当你感到失望时, 不要默默退缩到自己的内心世界, 而是勇敢地与伴侣沟通你的感受和需要。建议你对伴侣的缺点保持宽容, 完美的伴侣不存在, 但愿意一起成长的人就在你身边。在关系中, 你的深情和忠诚是宝贵的, 但也要学会在关系中保持自我。不要因为爱情而完全放弃自己的兴趣和追求, 一个有独立生活的你会让关系更加丰富多彩。记住, 真正的爱情不是两个完美的人在一起, 而是两个不完美的人选择一起面对世界的不完美。",
            'INTP': "在人际交往中, 你是一个独立、理性、有趣的人。你可能不是社交达人, 但对有趣的话题和深度的讨论充满热情。你最享受的社交方式是与志同道合的人进行智力交流 - 关于科学、哲学、技术或任何让你好奇的话题。你有独特的幽默感, 常常用出人意料的方式让人捧腹。你对人际关系的态度是质量胜过数量, 你宁愿有几个能进行深度对话的朋友, 也不愿有一大堆泛泛之交。与不同类型的人相处时, 理解情感型(F型)的人需要更多的情感回应和温暖表达, 而非只有逻辑分析。面对外向型(E型)的朋友时, 学会在社交活动中适度参与, 不要总是躲在角落里。与同样思考型的人在一起时, 你们能够进行令人兴奋的智力碰撞, 这是你最享受的社交体验。在亲密关系中, 你的独立性和独特思维是你的魅力所在。建议你也要学会表达情感和关怀, 不要让理性成为你唯一的表达方式。当伴侣分享感受时, 先给予情感上的支持和认可, 再进行分析和解决问题。学会计划一些浪漫和温馨的安排, 即使是小小的惊喜也会让你的伴侣感到被爱。你的智力和幽默是珍贵的, 但在关系中也要展现你柔软和温暖的一面。定期与伴侣分享你的内心世界, 让他们感受到你不仅是聪明的, 更是深情的。在日常相处中, 建议你放下手机和电脑, 给伴侣百分之百的关注。你可能会觉得同时做多件事很高效, 但在关系中, 全心的陪伴比高效更重要。当伴侣谈论他们的日常生活或感受时, 即使你不太感兴趣, 也要表现出真诚的倾听和关心。建议你尝试记住伴侣提到的重要日期和事件, 然后主动做出回应, 这种细节上的关注会让伴侣感到被珍视。在关系中, 你的智慧和幽默是吸引伴侣的重要因素, 但也要学会表达温柔和关怀。一个拥抱、一个关心的眼神、一句温暖的话, 这些看似简单的小事对你的伴侣来说可能比任何精妙的观点都更有意义。定期告诉伴侣你爱他们和欣赏他们的理由, 这将给关系注入持久的温暖。",
            'ESTP': "在人际交往中, 你是最有活力和魅力的存在。你走到哪里都能带来欢笑和活力, 是天生的社交达人。你喜欢与人互动, 享受刺激和冒险, 你的朋友们总是能在你身边找到乐趣。你有出色的观察力和应变能力, 在任何社交场合都能游刃有余。你的直率和坦诚让人感到轻松, 和你在一起不需要猜来猜去。你用实际的帮助和有趣的体验来表达关心。与不同类型的人相处时, 理解内向型(I型)的朋友需要更多的独处时间和安静的相处方式。面对直觉型(N型)的人时, 欣赏他们的深度思考和长远规划, 即使你觉得过于理论化。与同样外向型的人在一起时, 你们会碰撞出充满活力的火花, 但要注意不要过度刺激。在亲密关系中, 你的冒险精神和实际行动力是你的魅力。建议你也要学会享受安静的相处时光, 不是所有的浪漫都需要刺激和冒险。在关系中培养稳定性和持久性, 不要让追求新鲜感影响了你对承诺的坚守。当伴侣需要深度交流时, 放下手机, 认真倾听。学会表达更细腻的情感, 用言语和温柔的举动来补充你的行动表达。你的活力和热情是无可替代的, 用它来创造一段既刺激又稳固的亲密关系。在日常相处中, 建议你偶尔放慢节奏, 享受安静温馨的二人时光。不是每次约会都需要是刺激的冒险, 有时候一起做饭、散步或在沙发上看电影同样美好。当伴侣想要深入交流时, 放下手机, 给他们你完整的注意力。建议你对伴侣的情感需求保持敏感, 他们可能不会像你一样直接表达, 但同样需要被关心和被理解。在关系中, 你的活力和魅力让生活充满乐趣, 但也要学会在关系中展现稳定和承诺的一面。当你们经历困难时期时, 你的坚持和不离不弃将比任何刺激的约会都更让伴侣感动。记住, 最深沉的爱情不是轰轰烈烈的瞬间, 而是日复一日的陪伴和守护。",
            'ESFP': "在人际交往中, 你是最受欢迎和喜爱的人之一。你的热情、真诚和对生活的热爱让你成为任何聚会中的焦点。你有天生的社交能力, 能够让每个人都感到被关注和被欣赏。你善于创造快乐和美好的体验, 朋友们都很喜欢和你在一起。你的乐观精神和感染力能够点亮身边每个人的心情。你慷慨大方, 乐于分享你的时间、精力和资源。与不同类型的人相处时, 理解内向型(I型)的朋友需要更多的安静空间和一对一的交流。面对思考型(T型)的人时, 接受他们的直接和理性, 这不是在否定你。与同样情感型的人在一起时, 你们能够建立非常温暖和深厚的情感连接。在亲密关系中, 你通过创造美好体验和表达爱意来维系感情。建议你也要学会在日常生活中保持稳定和一致性, 浪漫不仅仅是惊喜和冒险。在关系中培养深度, 不要只停留在表面的快乐和轻松。当面对严肃的话题和深度的情感交流时, 不要逃避或转移话题。学会独处和自我反思, 这将帮助你更好地了解自己和伴侣的需求。你的温暖和快乐是珍贵的礼物, 用它来建立一段既充满欢笑又有深度的关系。在日常相处中, 建议你不仅用创造美好体验来表达爱, 也要学会用言语和倾听来表达关心。当伴侣有烦恼时, 安静地坐在他们身边倾听, 有时候比任何惊喜都更有力量。建议你对未来的生活有一些规划和打算, 即使你现在享受当下, 伴侣也可能需要一些安全感。在关系中, 你的热情和快乐是巨大的光芒, 但也要学会接纳和处理负面情绪。当你感到沮丧或不安时, 不要总是用笑容来掩饰, 真实地分享你的感受将让你们的关系更加真实和深入。记住, 一个健康的关系能够承载各种情绪, 不仅仅是快乐和欢笑。你的温暖和活力是珍贵的, 用它来创造一个让两个人都能真实做自己的关系空间。",
            'ENFP': "在人际交往中, 你是最有感染力和启发性的人。你的热情和创造力让你在任何社交场合都能脱颖而出。你善于发现每个人的潜力和闪光点, 这让身边的人在你面前感到被认可和鼓舞。你有广大的社交圈, 能够与各种各样的人建立联系。你对人性充满好奇, 喜欢了解每个人的故事。你的幽默、机智和真诚让你成为最受欢迎的朋友之一。与不同类型的人相处时, 理解实感型(S型)的人可能需要更多的具体细节和实际行动, 而非宏大的愿景。面对判断型(J型)的朋友时, 尽量遵守约定和时间, 这会让他们感到被尊重。与同样直觉型的人在一起时, 你们的对话会天马行空, 令人兴奋。在亲密关系中, 你的热情和创造力是巨大的优势。建议你也要学会保持专注和持久, 不要让新的可能性分散了你对当前关系的投入。在关系中培养稳定性和深度, 真正的爱情不仅是火花和激情, 更是日复一日的承诺和陪伴。当伴侣需要认真严肃的对话时, 放下你的幽默和跳脱, 给予他们你专注和真诚的回应。学会倾听而不急于给出解决方案, 有时候伴侣需要的只是被理解和被陪伴。你的热情和温暖是珍贵的, 用它来建立一段既充满激情又经得起时间考验的关系。在日常相处中, 建议你学会专注于眼前的人和当下的话题, 而不是总是展望未来或跳转到下一个有趣的想法。当伴侣与你分享重要的事情时, 给他们你完整的注意力和认真的回应。建议你对自己的承诺保持一致, 不要因为新的想法或机会而忽略了已经答应的事情。在关系中, 你的创意和热情让生活充满惊喜, 但也要学会在关系中创造稳定和可预测的日常仪式。固定的约会之夜、每周的深入对话、每天的一句感谢, 这些简单的习惯将给关系带来稳固的基础。记住, 爱情不仅需要激情和浪漫, 也需要坚持和守护。你的热情是珍贵的, 但有方向的热情才能创造真正持久的关系。",
            'ENTP': "在人际交往中, 你是最有趣和最具启发性的伙伴。你的机智、幽默和独特的思维角度让你在任何对话中都能带来新鲜的观点。你享受智力上的交锋和辩论, 这对你来说不是攻击, 而是一种交流方式。你有天生的魅力和说服力, 能够在社交场合中轻松吸引注意力。你对各种话题都感兴趣, 这使得你能够与不同领域的人进行有趣的对话。与不同类型的人相处时, 理解情感型(F型)的人可能在辩论中感到受伤, 学会区分讨论问题和攻击人格。面对实感型(S型)的人时, 尊重他们的务实观点, 不要总是用理论和可能性来挑战他们。与同样直觉型的人在一起时, 你们的对话会充满创意和启发, 但要注意有时也要脚踏实地。在亲密关系中, 你的聪明才智和幽默感是巨大的魅力。建议你也要学会表达真诚的情感和关怀, 不要总是用智慧和幽默来回避深层的情感交流。在关系中培养持久性和深度, 不要让对新鲜感的追求影响了你对伴侣的关注。当伴侣分享脆弱和感受时, 放下你的分析框架, 先给予温暖和理解。学会承诺和坚持, 让你的伴侣感受到你不仅仅是聪明的, 更是可靠的和深情的。定期与伴侣进行非辩论性的深度交流, 分享彼此的内心世界。在日常相处中, 建议你学会在辩论中放下胜负心, 有时候伴侣需要的不是一场精彩的辩论, 而是一个温暖的拥抱和一句我理解你。当伴侣分享他们的感受时, 不要立刻分析问题出在哪里, 而是先给予情感上的回应和支持。建议你对关系中的承诺保持认真和坚定的态度, 不要因为追求新鲜感而忽略了眼前人的感受。在关系中, 你的聪明和机智让人着迷, 但也要学会在适当的时候展现你的脆弱和依赖。告诉伴侣你需要他们, 这不是示弱, 而是让伴侣感受到自己被需要和被珍视。记住, 在亲密关系中, 最深刻的连接不是智力上的共鸣, 而是情感上的相互依赖和信任。你的智慧是珍贵的, 但在爱面前, 一颗温暖真诚的心同样重要。",
            'ESTJ': "在人际交往中, 你是一个可靠、有组织、直率的人。你不是最善言辞的人, 但你的行动力和可靠性让你成为朋友和同事中最值得信赖的人。你重视忠诚和责任, 认为这是人际关系的基石。你在社交场合中可能显得有些严肃, 但你的幽默感(虽然偏向冷幽默)其实很有魅力。你善于组织活动和协调事务, 是朋友圈中的大管家。与不同类型的人相处时, 理解情感型(F型)的人需要更多的情感回应和认可, 不要只关注问题和解决方案。面对感知型(P型)的朋友时, 尝试接受他们的灵活性和随性, 不是每个人都需要严格的时间表。与同样判断型的人在一起时, 你们能够高效协作, 共同完成目标。在亲密关系中, 你的可靠性和保护欲是你最大的魅力。建议你也要学会表达温柔和关怀, 不要让责任和义务成为你们之间唯一的纽带。在关系中关注伴侣的情感需求, 不只是提供物质和实际的保障。学会放松和享受生活, 偶尔放下工作和计划, 与伴侣享受轻松的时光。当伴侣表达情感时, 认真倾听和回应, 即使你觉得有些不理性。你的坚定和可靠是珍贵的, 用它来构建一段既有安全感又有温暖的关系。在日常相处中, 建议你偶尔放下计划和控制, 让事情自然地发展。惊喜和即兴的举动能给关系带来新鲜感和快乐。当伴侣想要倾诉时, 不要急于给出解决方案或评价, 先听他们把话说完, 表达你的理解和支持。建议你每周至少安排一次纯粹的约会时间, 不讨论工作、不处理家务, 只是享受彼此的陪伴。在关系中, 你的可靠和担当是巨大的安全感来源, 但也要学会表达柔软和浪漫。偶尔的一束花、一句我爱你、一个意外的约会, 这些小细节会让你的伴侣感受到被爱。记住, 最牢固的关系不是建立在责任和义务之上, 而是建立在爱、信任和相互欣赏之上。你的坚定是珍贵的, 但在爱情中, 温柔也是一种力量。",
            'ESFJ': "在人际交往中, 你是最温暖、最受欢迎的人之一。你天生善于建立和维护关系, 让每个人都感到被关心和被重视。你是聚会中的灵魂人物, 总是确保每个人都感到舒适和开心。你的善良、慷慨和体贴让你拥有广泛的社交圈和深厚的人际关系。你重视和谐, 善于调解矛盾, 让不同的人在同一空间中融洽相处。与不同类型的人相处时, 理解思考型(T型)的人可能更直接和理性, 这不是对你关心的否定。面对内向型(I型)的朋友时, 尊重他们对独处时间的需求, 不要过度热情。与同样情感型的人在一起时, 你们能够建立非常温暖和互助的关系。在亲密关系中, 你的关怀和付出是巨大的优势。建议你也要学会表达自己的需求和不满, 不要总是把伴侣的需求放在第一位。在关系中保持自己的独立性, 不要过度依赖伴侣的认可来获得自我价值。学会接受建设性的批评, 不要将其视为对你的否定。当面临冲突时, 不要总是回避或妥协, 有时正面解决问题才能让关系更加健康。你的温暖和关怀是珍贵的, 用它来建立一段既温馨又平等的关系, 让双方都能在其中成长和幸福。在日常相处中, 建议你学会接受伴侣的独立性和不同的处事方式, 不是每个人都像你一样善于照顾人, 伴侣的关心可能以你不熟悉的方式表达。当你感到被忽视或未被感激时, 直接而温和地与伴侣沟通, 而不是默默期待对方自己发现。建议你留出专门的时间来关注自己的需求和兴趣, 一个有自我关注的人才能在关系中给予更多。在关系中, 你的温暖和体贴是无价的, 但也要学会接受伴侣的批评和建议, 这不是对你付出的否定, 而是关系成长的必要部分。记住, 健康的关系中, 双方都有表达不满和需求的权利。你的关怀是珍贵的, 但不要让它成为束缚你自己的枷锁, 一个快乐满足的你才是最好的伴侣。",
            'ENFJ': "在人际交往中, 你是最有魅力和影响力的人之一。你有天生的感召力, 能够激励和鼓舞身边的每个人。你善于发现每个人的潜力和优点, 并帮助他们成为更好的自己。你是朋友圈中的精神领袖, 大家都信任你的判断和建议。你的同理心让你能够深刻理解每个人的处境, 你的表达能力让你能够将复杂的情感用简单的话语传达。你的人际关系以深度和意义为特点。与不同类型的人相处时, 理解思考型(T型)的人可能更需要逻辑推理而非情感支持。面对内向型(I型)的朋友时, 不要过度热情或试图拯救他们。与同样情感型的人在一起时, 你们能够建立非常深厚的情感纽带, 相互理解和支持。在亲密关系中, 你的理解和支持是无与伦比的。建议你也要学会倾听而非总是引导, 有时你的伴侣需要的只是被陪伴而非被指导。在关系中保持健康的边界感, 不要过度投入他人的需求而忽略了自己。学会接受伴侣的独立性, 不要试图控制或改变他们。当感到自己的需求被忽略时, 勇敢而温和地表达出来。你的影响力和关怀是珍贵的, 用它来建立一段相互成长、彼此成就的深度关系。在日常相处中, 建议你学会不要总是想要帮助和指导伴侣, 有时候他们只是需要一个平等的伴侣而非一个导师或拯救者。当你感到疲惫或需要支持时, 大胆地告诉伴侣, 你不需要总是坚强。建议你给伴侣足够的空间去做他们自己的事情, 包括犯错和从错误中学习。在关系中, 你的洞察力和关怀让你成为最理解伴侣的人, 但也要学会接受伴侣可能不理解你的所有感受。当你对关系的未来有担忧时, 与伴侣坦诚地交流, 而不是试图独自解决或调整。记住, 最深刻的爱不是试图把对方变成更好的版本, 而是接受和欣赏他们现在的样子。你的影响力是珍贵的, 但在亲密关系中, 平等和真实比任何感召力都更有力量。",
            'ENTJ': "在人际交往中, 你是一个自信、有魄力、有远见的人。你的领导气质和决断力让你在社交场合中自然而然地成为焦点。你不善于闲聊, 但对有深度的话题和有挑战性的讨论充满热情。你的人际关系以尊重和共同目标为基础, 你尊重有能力、有担当的人。你的坦率和直爽让有些人觉得你难以接近, 但真正了解你的人知道你有非常重情义的一面。与不同类型的人相处时, 理解情感型(F型)的人需要更多的温暖和情感回应, 不要只关注效率和结果。面对感知型(P型)的人时, 尝试接受他们的灵活性和随性, 并非每个人都需要严格的计划和结构。与同样外向型的人在一起时, 你们可以成为强大的合作伙伴, 但要注意分配领导权。在亲密关系中, 你的保护和担当是巨大的魅力。建议你也要学会展现脆弱和温柔, 不要让坚强成为你唯一的面具。在关系中关注伴侣的感受和需求, 不只是提供物质和实际的保障。学会放慢脚步, 享受关系中的平凡时刻。当伴侣表达情感时, 认真倾听而不是急于分析和解决。你的自信和远见是珍贵的, 用它来构建一段既有方向又有温度的关系, 让伴侣感受到在你身边既有安全感又有成长空间。在日常相处中, 建议你放下控制欲, 让伴侣有自己做决定和犯错的空间。在你的世界里, 高效和正确是最重要的, 但在关系中, 感受和过程同样重要。当伴侣分享他们的成就时, 给予真诚的赞美和欣赏, 而不是指出还可以做得更好的地方。建议你定期安排一些不涉及任何目标的纯粹约会, 享受与伴侣在一起的时间本身。在关系中, 你的远见和魄力让生活充满方向感, 但也要学会在伴侣面前展现你的温柔和脆弱。当你感到压力或困惑时, 向伴侣倾诉, 这不是失去力量, 而是加深你们之间信任的方式。记住, 最强大的领导者也是最懂得关爱的人。你的自信是珍贵的, 但在爱情中, 一颗柔软真诚的心同样重要。",
        }
        return advice.get(mbti_type, '')


    def _calculate_duration(self, start_time, end_time):
        if not start_time or not end_time:
            return None
        try:
            start = datetime.fromisoformat(start_time.replace('Z', '+00:00'))
            end = datetime.fromisoformat(end_time.replace('Z', '+00:00'))
            duration = end - start
            total_seconds = duration.total_seconds()
            minutes = int(total_seconds // 60)
            seconds = int(total_seconds % 60)
            return f"{minutes}分{seconds}秒"
        except:
            return None

    def _get_analytics_overview(self):
        """概览统计"""
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 总用户数
            cursor.execute("SELECT COUNT(*) FROM users")
            total_users = cursor.fetchone()[0]
            
            # 总测试次数
            cursor.execute("SELECT COUNT(*) FROM test_records")
            total_tests = cursor.fetchone()[0]
            
            # 已完成测试次数
            cursor.execute("SELECT COUNT(*) FROM test_records WHERE status = 'completed'")
            completed_tests = cursor.fetchone()[0]
            
            # 平均完成时间（秒）
            cursor.execute("""
                SELECT AVG(
                    (julianday(end_time) - julianday(start_time)) * 86400
                ) FROM test_records 
                WHERE status = 'completed' AND end_time IS NOT NULL
            """)
            avg_duration = cursor.fetchone()[0]
            avg_duration_str = f"{int(avg_duration // 60)}分{int(avg_duration % 60)}秒" if avg_duration else "0 分 0 秒"
            
            # 今日新增用户
            cursor.execute("""
                SELECT COUNT(*) FROM users 
                WHERE DATE(created_at) = DATE('now', 'localtime')
            """)
            today_new_users = cursor.fetchone()[0]
            
            # 今日测试次数
            cursor.execute("""
                SELECT COUNT(*) FROM test_records 
                WHERE DATE(start_time) = DATE('now', 'localtime')
            """)
            today_tests = cursor.fetchone()[0]
            
            self._send_success({
                'total_users': total_users,
                'total_tests': total_tests,
                'completed_tests': completed_tests,
                'avg_completion_time': avg_duration_str,
                'today_new_users': today_new_users,
                'today_tests': today_tests
            })
        except Exception as e:
            self._send_error(500, f"获取概览统计失败：{str(e)}")
        finally:
            conn.close()
    
    def _get_mbti_distribution(self):
        """MBTI 类型分布"""
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            parsed_path = urlparse(self.path)
            query_params = parse_qs(parsed_path.query)
            
            start_date = query_params.get('start_date', [None])[0]
            end_date = query_params.get('end_date', [None])[0]
            
            # 构建查询条件
            where_clause = "WHERE status = 'completed' AND mbti_type IS NOT NULL"
            if start_date:
                where_clause += f" AND DATE(end_time) >= DATE('{start_date}')"
            if end_date:
                where_clause += f" AND DATE(end_time) <= DATE('{end_date}')"
            
            # 获取总数
            cursor.execute(f"SELECT COUNT(*) FROM test_records {where_clause}")
            total = cursor.fetchone()[0]
            
            # 获取各类型分布
            cursor.execute(f"""
                SELECT mbti_type, COUNT(*) as count 
                FROM test_records 
                {where_clause}
                GROUP BY mbti_type 
                ORDER BY count DESC
            """)
            
            distribution = []
            mbti_types = ['ISTJ', 'ISFJ', 'INFJ', 'INTJ', 'ISTP', 'ISFP', 'INFP', 'INTP',
                         'ESTP', 'ESFP', 'ENFP', 'ENTP', 'ESTJ', 'ESFJ', 'ENFJ', 'ENTJ']
            
            type_counts = {row[0]: row[1] for row in cursor.fetchall()}
            
            for mbti_type in mbti_types:
                count = type_counts.get(mbti_type, 0)
                percentage = round(count / total * 100, 2) if total > 0 else 0
                distribution.append({
                    'type': mbti_type,
                    'count': count,
                    'percentage': percentage
                })
            
            self._send_success({
                'total': total,
                'distribution': distribution,
                'filter': {
                    'start_date': start_date,
                    'end_date': end_date
                }
            })
        except Exception as e:
            self._send_error(500, f"获取 MBTI 分布失败：{str(e)}")
        finally:
            conn.close()
    
    def _get_dimension_distribution(self):
        """四维度分布"""
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 计算各维度分数总和
            cursor.execute("""
                SELECT 
                    SUM(score_e) as total_e,
                    SUM(score_i) as total_i,
                    SUM(score_s) as total_s,
                    SUM(score_n) as total_n,
                    SUM(score_t) as total_t,
                    SUM(score_f) as total_f,
                    SUM(score_j) as total_j,
                    SUM(score_p) as total_p
                FROM answers
            """)
            
            row = cursor.fetchone()
            total_e = row[0] or 0
            total_i = row[1] or 0
            total_s = row[2] or 0
            total_n = row[3] or 0
            total_t = row[4] or 0
            total_f = row[5] or 0
            total_j = row[6] or 0
            total_p = row[7] or 0
            
            # 计算各维度总数
            ei_total = total_e + total_i
            sn_total = total_s + total_n
            tf_total = total_t + total_f
            jp_total = total_j + total_p
            
            distribution = {
                'E_I': {
                    'E': {'count': total_e, 'percentage': round(total_e / ei_total * 100, 2) if ei_total > 0 else 0},
                    'I': {'count': total_i, 'percentage': round(total_i / ei_total * 100, 2) if ei_total > 0 else 0}
                },
                'S_N': {
                    'S': {'count': total_s, 'percentage': round(total_s / sn_total * 100, 2) if sn_total > 0 else 0},
                    'N': {'count': total_n, 'percentage': round(total_n / sn_total * 100, 2) if sn_total > 0 else 0}
                },
                'T_F': {
                    'T': {'count': total_t, 'percentage': round(total_t / tf_total * 100, 2) if tf_total > 0 else 0},
                    'F': {'count': total_f, 'percentage': round(total_f / tf_total * 100, 2) if tf_total > 0 else 0}
                },
                'J_P': {
                    'J': {'count': total_j, 'percentage': round(total_j / jp_total * 100, 2) if jp_total > 0 else 0},
                    'P': {'count': total_p, 'percentage': round(total_p / jp_total * 100, 2) if jp_total > 0 else 0}
                }
            }
            
            self._send_success({'distribution': distribution})
        except Exception as e:
            self._send_error(500, f"获取维度分布失败：{str(e)}")
        finally:
            conn.close()
    
    def _get_question_stats(self):
        """题目分析"""
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 获取每道题的选项分布
            cursor.execute("""
                SELECT question_id, answer, COUNT(*) as count 
                FROM answers 
                GROUP BY question_id, answer 
                ORDER BY question_id, count DESC
            """)
            
            question_stats = {}
            for row in cursor.fetchall():
                question_id = row[0]
                answer = row[1]
                count = row[2]
                
                if question_id not in question_stats:
                    question_stats[question_id] = {'question_id': question_id, 'options': []}
                
                question_stats[question_id]['options'].append({
                    'answer': answer,
                    'count': count
                })
            
            # 添加题目文本
            for q_id in question_stats:
                for q in MBTI_QUESTIONS:
                    if q['id'] == q_id:
                        question_stats[q_id]['question'] = q['question']
                        break
            
            self._send_success({
                'total_questions': len(question_stats),
                'stats': list(question_stats.values())
            })
        except Exception as e:
            self._send_error(500, f"获取题目统计失败：{str(e)}")
        finally:
            conn.close()
    
    def _get_analytics_trends(self):
        """时间趋势"""
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            parsed_path = urlparse(self.path)
            query_params = parse_qs(parsed_path.query)
            
            days = int(query_params.get('days', ['7'])[0])
            
            # 获取每日趋势
            cursor.execute(f"""
                SELECT DATE(start_time) as date, COUNT(*) as count 
                FROM test_records 
                WHERE DATE(start_time) >= DATE('now', 'localtime', '-{days} days')
                GROUP BY DATE(start_time)
                ORDER BY date
            """)
            
            test_trends = [{'date': row[0], 'count': row[1]} for row in cursor.fetchall()]
            
            # 获取每日注册趋势
            cursor.execute(f"""
                SELECT DATE(created_at) as date, COUNT(*) as count 
                FROM users 
                WHERE DATE(created_at) >= DATE('now', 'localtime', '-{days} days')
                GROUP BY DATE(created_at)
                ORDER BY date
            """)
            
            register_trends = [{'date': row[0], 'count': row[1]} for row in cursor.fetchall()]
            
            # 获取每日完成趋势
            cursor.execute(f"""
                SELECT DATE(end_time) as date, COUNT(*) as count 
                FROM test_records 
                WHERE status = 'completed' 
                AND DATE(end_time) >= DATE('now', 'localtime', '-{days} days')
                GROUP BY DATE(end_time)
                ORDER BY date
            """)
            
            completion_trends = [{'date': row[0], 'count': row[1]} for row in cursor.fetchall()]
            
            self._send_success({
                'days': days,
                'register_trends': register_trends,
                'test_trends': test_trends,
                'completion_trends': completion_trends
            })
        except Exception as e:
            self._send_error(500, f"获取趋势数据失败：{str(e)}")
        finally:
            conn.close()
    
    def _get_completion_rate(self):
        """完成率分析"""
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 总测试数
            cursor.execute("SELECT COUNT(*) FROM test_records")
            total_tests = cursor.fetchone()[0] or 0
            
            # 完成测试数
            cursor.execute("SELECT COUNT(*) FROM test_records WHERE status = 'completed'")
            completed_tests = cursor.fetchone()[0] or 0
            
            # 未完成测试数
            in_progress_tests = total_tests - completed_tests
            
            # 完成率
            completion_rate = round(completed_tests / total_tests * 100, 2) if total_tests > 0 else 0
            
            # 未完成测试的平均答题数
            cursor.execute("""
                SELECT AVG(answer_count) FROM (
                    SELECT test_record_id, COUNT(*) as answer_count 
                    FROM answers 
                    GROUP BY test_record_id
                ) WHERE test_record_id IN (
                    SELECT id FROM test_records WHERE status = 'in_progress'
                )
            """)
            avg_answers_incomplete = cursor.fetchone()[0] or 0
            
            # 放弃率（未完成的比例）
            abandonment_rate = round(in_progress_tests / total_tests * 100, 2) if total_tests > 0 else 0
            
            self._send_success({
                'total_tests': total_tests,
                'completed_tests': completed_tests,
                'in_progress_tests': in_progress_tests,
                'completion_rate': completion_rate,
                'avg_answers_incomplete': round(avg_answers_incomplete, 2),
                'abandonment_rate': abandonment_rate
            })
        except Exception as e:
            self._send_error(500, f"获取完成率分析失败：{str(e)}")
        finally:
            conn.close()
    
    def _get_demographics(self):
        """用户画像"""
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()
        
        try:
            # 年龄分布
            cursor.execute("""
                SELECT 
                    CASE 
                        WHEN age < 18 THEN '18 岁以下'
                        WHEN age BETWEEN 18 AND 24 THEN '18-24 岁'
                        WHEN age BETWEEN 25 AND 34 THEN '25-34 岁'
                        WHEN age BETWEEN 35 AND 44 THEN '35-44 岁'
                        WHEN age BETWEEN 45 AND 54 THEN '45-54 岁'
                        WHEN age BETWEEN 55 AND 64 THEN '55-64 岁'
                        ELSE '65 岁以上'
                    END as age_group,
                    COUNT(*) as count
                FROM users
                GROUP BY age_group
                ORDER BY age_group
            """)
            
            age_distribution = [{'age_group': row[0], 'count': row[1]} for row in cursor.fetchall()]
            
            # 性别分布
            cursor.execute("""
                SELECT gender, COUNT(*) as count 
                FROM users 
                GROUP BY gender
            """)
            
            total_users = sum(row[1] for row in cursor.fetchall())
            
            # 重新查询以获取详细数据
            cursor.execute("""
                SELECT gender, COUNT(*) as count 
                FROM users 
                GROUP BY gender
            """)
            
            gender_distribution = []
            for row in cursor.fetchall():
                gender_distribution.append({
                    'gender': row[0],
                    'count': row[1],
                    'percentage': round(row[1] / total_users * 100, 2) if total_users > 0 else 0
                })
            
            self._send_success({
                'age_distribution': age_distribution,
                'gender_distribution': gender_distribution,
                'total_users': total_users
            })
        except Exception as e:
            self._send_error(500, f"获取用户画像失败：{str(e)}")
        finally:
            conn.close()

def init_database():
    """初始化数据库"""
    conn = sqlite3.connect(DB_FILE)
    cursor = conn.cursor()
    
    # 创建用户表
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            nickname TEXT NOT NULL,
            gender TEXT NOT NULL,
            age INTEGER NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    ''')
    
    # 创建测试记录表
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS test_records (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            test_id TEXT NOT NULL,
            mbti_type TEXT,
            status TEXT DEFAULT 'in_progress',
            start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            end_time TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (id)
        )
    ''')
    
    # 创建答案详情表
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS answers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            test_record_id INTEGER NOT NULL,
            question_id INTEGER NOT NULL,
            answer TEXT NOT NULL,
            score_e INTEGER DEFAULT 0,
            score_i INTEGER DEFAULT 0,
            score_s INTEGER DEFAULT 0,
            score_n INTEGER DEFAULT 0,
            score_t INTEGER DEFAULT 0,
            score_f INTEGER DEFAULT 0,
            score_j INTEGER DEFAULT 0,
            score_p INTEGER DEFAULT 0,
            FOREIGN KEY (test_record_id) REFERENCES test_records (id)
        )
    ''')
    
    conn.commit()
    conn.close()
    print("数据库表创建完成")

def main():
    """主函数"""
    print(f"正在启动MBTI测试服务器...")
    print(f"端口: {PORT}")
    print(f"数据库: {DB_FILE}")
    
    # 初始化数据库
    init_database()
    
    print("数据库初始化完成")
    print(f"题目数量: {len(MBTI_QUESTIONS)}")
    print("服务器正在启动...")
    
    try:
        httpd = HTTPServer(('0.0.0.0', PORT), MBTIServer)
        print(f"服务器已启动，访问 http://localhost:{PORT}")
        print("按 Ctrl+C 停止服务器")
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\n服务器已停止")
    except Exception as e:
        print(f"服务器启动失败: {e}")

if __name__ == '__main__':
    main()
