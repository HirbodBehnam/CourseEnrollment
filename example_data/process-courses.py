import json
import jdatetime

# Load json
f = open('courses.json')
data = json.load(f)['message']
f.close()

# Check each course
f = open('courses.sql', 'w')
for course in data:
	courseID = course['number']
	groupID = course['group']
	units = course['units']
	department = course['department']
	capacity = course['capacity']
	reserve = 5 # no data from backend
	lecturer = course['instructors']
	title = course['title']
	notes = course['description']
	sex_lock = "NULL"
	if 'خانم' in notes or 'خواهران' in notes:
		sex_lock = "'female'"
	elif 'آقایان' in notes or 'برادران' in notes:
		sex_lock = "'male'"
	if len(course['schedule']) == 0:
		print('course with empty schedule:', course['id'])
		continue
	classTime = int(course['schedule'][0]['start'] * 60) + 1440 * int(course['schedule'][0]['end'] * 60)
	for day in course['schedule']:
		classTime = classTime | (1 << (day['day'] + 21))
	examTime = 'NULL'
	if course['examDate'] != " ":
		try:
			examTime = "'" + jdatetime.datetime.strptime(course['examDate'], "%Y/%m/%d %H:%M").togregorian().strftime("%Y/%m/%d %H:%M:%S") + "'"
		except:
			pass
	
	f.write(f"INSERT INTO courses (course_id, group_id, for_department, name, lecturer, units, capacity, reserve_capacity, exam_time, class_time, sex_lock, notes) VALUES ({courseID}, {groupID}, {department}, '{title}', '{lecturer}', {units}, {capacity}, {reserve}, {examTime}, {classTime}, {sex_lock}, '{notes}');\n")
	
f.close()