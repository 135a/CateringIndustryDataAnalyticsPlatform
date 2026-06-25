import os
import time
import random
import pandas as pd
import pymysql
from datetime import datetime

# ================= 配置区域 =================
# MySQL 数据库配置 (优先读取环境变量，适配 Docker)
DB_CONFIG = {
    'host': os.getenv('DB_HOST', '127.0.0.1'),
    'port': 3306,
    'user': 'root',
    'password': 'root',  # 请在运行前将其替换为您本地的MySQL真实密码
    'database': 'catering_data',
    'charset': 'utf8mb4'
}

# Excel 保存目录 (系统会自动创建该目录)
EXCEL_SAVE_DIR = os.getenv('EXCEL_DIR', r"D:\CateringIndustryDataAnalyticsPlatform\data\excel")
os.makedirs(EXCEL_SAVE_DIR, exist_ok=True)

# 目标城市与行政区示例
CITY = '杭州市'
DISTRICTS = ['西湖区', '上城区', '拱墅区', '滨江区', '萧山区', '余杭区', '临平区', '钱塘区', '富阳区', '临安区']
CATEGORIES = [
    '小吃快餐', '咖啡', '自助餐', '面包甜点', '酒吧', 
    '烧烤烤串', '创意菜', '鱼鲜海鲜', '水果生鲜', '饮品',
    '特色菜', '地方菜系', '食品滋补', '农家菜', '私房菜',
    '家常菜', '粤菜', '川菜', '面馆', '江浙菜',
    '火锅', '茶馆', '西餐', '日式料理', '小龙虾',
    '烤肉', '湘菜', '东北菜', '北京菜', '韩式料理'
]

# 抓取数据量控制：每个分类下每次抓取的商户数量区间 (您可以随时调大此数值以获取更多数据)
FETCH_COUNT_MIN = 5
FETCH_COUNT_MAX = 12

# ============================================

def init_db_connection():
    """初始化数据库连接（增加重试机制以等待 Docker 中的 MySQL 启动完成）"""
    max_retries = 10
    for i in range(max_retries):
        try:
            conn = pymysql.connect(**DB_CONFIG)
            print("   [+] 成功连接到 MySQL 数据库！")
            return conn
        except Exception as e:
            print(f"⚠️ 数据库连接失败，等待 MySQL 启动... ({i+1}/{max_retries}): {e}")
            time.sleep(5)
    return None

def fetch_data_from_dianping(district, category):
    """
    模拟从大众点评爬取数据。
    注意：大众点评有极强的反爬机制（如：字体反爬、CSS偏移、封IP）。
    实际开发中，这里通常需要接入 Playwright/Selenium 或逆向破解其接口，并配合代理IP池和Cookie池。
    为演示“双写存储”流程，此处随机生成结构化餐饮数据。
    """
    print(f"[*] 正在抓取: {CITY} - {district} - {category} 的数据...")
    
    mock_data = []
    for i in range(1, random.randint(FETCH_COUNT_MIN, FETCH_COUNT_MAX)): # 根据配置控制每次抓取的商户数量
        shop_id = f"dp_{int(time.time())}_{random.randint(1000, 9999)}"
        mock_data.append({
            'shop_id': shop_id,
            'name': f"{district}正宗{category}_{i}店",
            'category_name': category,
            'district_name': district,
            'address': f"杭州市{district}某某街道{random.randint(1, 999)}号",
            'avg_price': round(random.uniform(20, 400), 2),
            'rating': round(random.uniform(3.5, 5.0), 1),
            'review_count': random.randint(10, 5000),
            'opening_hours': '10:00-22:00',
            'taste_score': round(random.uniform(3.5, 5.0), 1),
            'environment_score': round(random.uniform(3.5, 5.0), 1),
            'service_score': round(random.uniform(3.5, 5.0), 1),
            'has_free_parking': random.choice([0, 1]),
            'is_reservable': random.choice([0, 1]),
            'has_baby_chair': random.choice([0, 1]),
            'has_private_room': random.choice([0, 1]),
        })
    return mock_data

def save_to_mysql(conn, data_list):
    """将抓取到的数据写入 MySQL 数据库"""
    if not data_list or not conn:
        return
        
    cursor = conn.cursor()
    try:
        # 首先确保存储数据所需的分类和行政区已经存在于字典表中，避免插入时因找不到外键而报错
        for data in data_list:
            cursor.execute("INSERT IGNORE INTO categories (name) VALUES (%s)", (data['category_name'],))
            cursor.execute("INSERT IGNORE INTO districts (city_name, district_name) VALUES (%s, %s)", (CITY, data['district_name']))
        conn.commit()

        # 此处利用子查询自动将名称转为字典表的 category_id 和 district_id
        # ON DUPLICATE KEY UPDATE 用于应对商户数据更新（如果商户已存在，则更新最新评分和评价数）
        insert_sql = """
            INSERT INTO restaurants (
                shop_id, name, category_id, district_id, address, 
                avg_price, rating, review_count, 
                opening_hours, taste_score, environment_score, service_score,
                has_free_parking, is_reservable, has_baby_chair, has_private_room
            ) VALUES (
                %(shop_id)s, %(name)s, 
                (SELECT id FROM categories WHERE name = %(category_name)s LIMIT 1), 
                (SELECT id FROM districts WHERE district_name = %(district_name)s LIMIT 1), 
                %(address)s, %(avg_price)s, %(rating)s, 
                %(review_count)s, %(opening_hours)s, %(taste_score)s, %(environment_score)s, %(service_score)s,
                %(has_free_parking)s, %(is_reservable)s, %(has_baby_chair)s, %(has_private_room)s
            )
            ON DUPLICATE KEY UPDATE 
                rating = VALUES(rating), 
                review_count = VALUES(review_count),
                avg_price = VALUES(avg_price),
                has_free_parking = VALUES(has_free_parking),
                is_reservable = VALUES(is_reservable),
                has_baby_chair = VALUES(has_baby_chair),
                has_private_room = VALUES(has_private_room);
        """
        cursor.executemany(insert_sql, data_list)
        conn.commit()
            
        print(f"   [+] 成功将 {len(data_list)} 条商户数据存入 MySQL 数据库。")
    except Exception as e:
        conn.rollback()
        print(f"   [-] 写入 MySQL 失败: {e}")
    finally:
        cursor.close()

def save_to_excel(data_list):
    """将数据导出到 Excel 文件作为本地双写备份"""
    if not data_list:
        return
        
    df = pd.DataFrame(data_list)
    
    # 按照当前日期生成 Excel 文件名，例如：dianping_data_20231025.xlsx
    date_str = datetime.now().strftime("%Y%m%d")
    filename = f"dianping_data_{date_str}.xlsx"
    filepath = os.path.join(EXCEL_SAVE_DIR, filename)
    
    # 双写逻辑：如果今日文件已存在，则追加内容；如果不存在，则创建新文件。
    if os.path.exists(filepath):
        existing_df = pd.read_excel(filepath)
        updated_df = pd.concat([existing_df, df], ignore_index=True)
        # 根据第三方平台ID(shop_id)进行去重，保留最新数据
        updated_df.drop_duplicates(subset=['shop_id'], keep='last', inplace=True)
        updated_df.to_excel(filepath, index=False)
        print(f"   [+] 追加写入并更新本地 Excel 备份: {filepath}")
    else:
        df.to_excel(filepath, index=False)
        print(f"   [+] 新建本地 Excel 备份: {filepath}")

def run_spider():
    print("="*60)
    print("🚀 启动大众点评餐饮数据爬虫程序...")
    print(f"📁 本地 Excel 存储目录将被设置在:\n   {EXCEL_SAVE_DIR}")
    print("="*60)
    
    conn = init_db_connection()
    if not conn:
        print("💡 提示：目前未连接到数据库，但爬虫程序将继续运行，抓取的数据将仅保存在 Excel 中。")
    
    all_data = [] # 用于最终统一导出Excel的容器
    
    # 遍历各大区与分类进行爬取
    for dist in DISTRICTS:
        for cat in CATEGORIES:
            data = fetch_data_from_dianping(dist, cat)
            all_data.extend(data)
            
            # 获取到一批数据，就存一次数据库，防止数据丢失
            if conn:
                save_to_mysql(conn, data)
                
            # 取消随机休眠以实现极速生成测试数据
            # time.sleep(random.uniform(1.5, 3.5)) 
            
    print("\n📦 全部分类抓取完毕，开始生成本地 Excel 备份...")
    # 一次性将今日所有抓取到的数据备份进 Excel
    save_to_excel(all_data)
    
    if conn:
        conn.close()
        
    print("✅ 今日爬虫任务执行完毕！")

if __name__ == "__main__":
    # 需要提前安装依赖包：pip install pandas pymysql openpyxl
    run_spider()
