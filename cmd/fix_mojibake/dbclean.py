from ftfy import fix_and_explain
import sqlite3

conn = sqlite3.connect("dev.db")
rows = conn.execute("SELECT id, title FROM items").fetchall()

fixes = {}

print("searching for mojibake...")
for item in rows:
    fixed, explanation = fix_and_explain(item[1])
    if len(explanation) > 0:
        print(f"#{item[0]}\t'{item[1]}' --> '{fixed}'")
        fixes[item[0]] = fixed

print(f"Found {len(fixes)} broken titles. Updating database with fixes...")
for id in fixes:
    print(f"{id}: {fixes[id]}")
    conn.execute("UPDATE items SET title=? WHERE id=?", (fixes[id], id))

conn.commit()

print("done")

