path = r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views\ModalTest.vue"
with open(path, 'rb') as f:
    data = f.read()
text = data.decode('utf-8', errors='replace')
idx = text.find("<style scoped>")
print(f'<style scoped> at: {idx}')
print(f'Before: {repr(text[idx-50:idx+20])}')
