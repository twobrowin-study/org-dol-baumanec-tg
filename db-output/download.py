from io import BytesIO
import pandas as pd
from sqlalchemy import create_engine
from flask import Flask, send_file
from datetime import datetime
from os import environ

app = Flask(__name__)

connection = create_engine(environ['CONNECTION_STRING'])

@app.route('/download')
def download():
    df = pd.read_sql("select * from user_outdata", connection)
    df.columns=["Имя", "Фамилия", "Пол", "Специалист", "Пройден"]
    df['Пол'] = df['Пол'].apply(lambda x: 'Муж.' if x == 'male' else 'Жен.')
    df['Пройден'] = df['Пройден'].apply(lambda x: 'Да' if x == 1 else 'Нет')
    
    output = BytesIO()
    writer = pd.ExcelWriter(output, engine='xlsxwriter')

    sheetname = "Прохождение медосмотра"
    df.to_excel(writer, startrow=0, merge_cells=False, sheet_name=sheetname, index=False)

    worksheet = writer.sheets[sheetname]
    row_count = len(df.index)
    column_count = len(df.columns)
    worksheet.autofilter(0, 0, row_count-1, column_count-1)  
    for idx, col in enumerate(df):
        series = df[col]
        max_len = max((
            series.astype(str).map(len).max(),
            len(str(series.name))
            )) + 5
        worksheet.set_column(idx, idx, max_len) 
    writer.save()

    writer.close()
    output.seek(0)
        
    now_datetime = datetime.now().strftime("%d-%m-%Y_%H-%M-%S")
    filename = "med_{}.xlsx".format(now_datetime)

    return send_file(output, attachment_filename=filename, as_attachment=True)

def create_app():
   return app

if __name__ == '__main__':
   app.run(debug=True)