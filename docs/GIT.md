# Правила работы с гитом

- Для каждой задачи мы создаем отдельную ветку (фича ветку).
- Фича ветки называем "IVI-{номер задачи}-{описание для себя}".
- Все коммиты в ветках должны начинаться с номера задачи.
  - IVI-0: fix some issues
- При мердже фича ветки в main коммиты сквошатся.
- Названия всех коммитов в main ветке должны совпадать с названием задачи, в рамках которой был создан коммит.

## Реквесты
 
Названия реквестов должны полностью совпадать с названием задачи в YouTrack.

Например, задача нызвается "IVI-12 Написать в репозиториях правила работы с гитом", тогда название реквеста - "IVI-12 Написать в репозиториях правила работы с гитом".


## Git Hooks

Для удобства работы с гитом были созданы pre-push и prepare-commit-msg хуки.

- pre-push проверяет имя ветки на корректность;
- prepare-commit-msg добавляет к коммиту номер задачи.

### Установка
Положить представленные ниже хуки в папку .git/hooks в своем репозитории. Сделать файлы с хуками исполняемыми.

### Хуки

#### pre-push
```python
#!/usr/bin/python3
import re

def main():
    branch_name = get_current_branch_name()
    if not re.match(r'IVI-\d+', branch_name):
        error_description = f'invalid branch name: {branch_name}'
        raise Exception(error_description)

def get_current_branch_name():
    with open('.git/HEAD', 'r') as f:
        ref = f.readline().strip()
    if ref.startswith('ref:'):
        branch_name = ref[16:]
        return branch_name

if __name__ == "__main__":
    main()

```

#### prepare-commit-msg

```python
#!/usr/bin/python3
import re
import sys

def main():
    branch_name = get_current_branch_name()
    if branch_name:
        task_number = extract_task_number(branch_name)
        if task_number:
            prepend_task_number(task_number)

def get_current_branch_name():
    with open('.git/HEAD', 'r') as f:
        ref = f.readline().strip()
    if ref.startswith('ref:'):
        branch_name = ref[16:]
        return branch_name

def extract_task_number(branch_name):
    match = re.search(r'IVI-\d+', branch_name)  
    if match:
        return match.group(0)

def prepend_task_number(task_number):
    commit_msg_file = sys.argv[1]
    with open(commit_msg_file, 'r+') as f:
        commit_msg = f.read()
        new_msg = f'{task_number}: {commit_msg}'
        f.seek(0, 0)
        f.write(new_msg)

if __name__ == "__main__":
    main()

```
