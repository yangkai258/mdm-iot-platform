path = r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views\ai\AIMonitorView.vue"
with open(path, 'rb') as f:
    data = f.read()

# Find ALL variations of <a-statistic
print('<a-statistic: ', data.count(b'<a-statistic'))
print('</a-statistic>: ', data.count(b'</a-statistic>'))

print('<a-card: ', data.count(b'<a-card'))
print('</a-card>: ', data.count(b'</a-card>'))

# Find position
idx = data.find(b'</a-statistic>')
print(f'First </a-statistic> at: {idx}')
print(f'Context: {data[idx-50:idx+30]}')