import psycopg2
import uuid
import json

with open('HTN_2023_BE_Challenge_Data.json', 'r') as file:
    data = json.load(file)

def etl_user_skills(conn, data):
    with conn.cursor() as cur:
        for user in data:
            print(f"Processing user: {user['name']}")
            user_id = str(uuid.uuid4())
            cur.execute("""
                        INSERT INTO users (id, name, company, email, phone) VALUES (%s, %s, %s, %s, %s) ON CONFLICT (email) DO NOTHING;
                        """, (user_id, user['name'], user['company'], user['email'], user['phone']))
            cur.execute("""
                        SELECT id FROM users WHERE email = %s;
                        """, (user['email'],))
            user_id = cur.fetchone()[0]
            
            for skill in user['skills']:
                skill_id = str(uuid.uuid4())
                cur.execute("""
                            INSERT INTO skills (id, name) VALUES (%s, %s) ON CONFLICT (name) DO NOTHING;
                            """, (skill_id, skill['skill']))
                cur.execute("""
                            SELECT id FROM skills WHERE name = %s;
                            """, (skill['skill'],))
                skill_id = cur.fetchone()[0]
                
                # Insert the user-skill relationship
                cur.execute("""
                            INSERT INTO users_skills (user_id, skill_id, rating) VALUES (%s, %s, %s) ON CONFLICT (user_id, skill_id) DO UPDATE SET rating = EXCLUDED.rating;
                            """, (user_id, skill_id, skill['rating']))
                
        conn.commit()


connection_config = {
    "dbname": "postgres",
    "user": "user",
    "password": "password",
    "host": "localhost"
}

if __name__ == "__main__":
    conn = psycopg2.connect(**connection_config)
    try:
        print("Executing ETL for user and skills...")
        etl_user_skills(conn, data)
    finally:
        conn.close()