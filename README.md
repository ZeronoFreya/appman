# appman

* 检查配置中定义的软链接，如果没有映射链接 ( 包括链接是否正确 )，先映射链接后再执行程序
* 检查配置中定义的需要管理员权限的程序，以静默方式 ( 无UAC提示 ) 启动
* 如果传入的程序没有记录，则创建快捷方式 (**_as_appman.lnk) 并写入配置文件，请自行修改 `config.json`


### 第一次使用或重做系统后务必运行`安装.bat` ( 管理员权限 )，功能是添加必要的计划任务

*   noUAC.runadmin
*   noUAC.mksymlink

> 通过计划任务启动的程序，不会提示UAC，且它调用的程序会继承管理员权限

### 为什么需要这个程序

第一个功能: 
以 `r2modman.exe` 为例 (见最后的`config.json`)，这是一个mod管理器，但无法指定安装目录及关闭更新，它会自动安装在 `%localappdata%\\Programs\\r2modman` 及 下载更新到 `%localappdata%\\r2modman-updater` ，无论你把 `r2modman` 移植到哪里，它都会检测更新并安装在C盘。
解决办法就是创建软链接
> 用来修改各种游戏的存档也是个不错的选择

第二个功能就很纯粹了，windows没有提供UAC白名单，不想关UAC也不想要UAC提示

---

## 使用方法

### 手动

*   准备一个快捷方式，如右键appman.exe，创建快捷方式
*   右键打开属性
*   目标栏改为: `E:\...\appman.exe "C:\Windows\notepad.exe" "d:\demo.txt"`
    *   用notepad打开demo.txt
*   其他项随意，比如可以改一下图标

### 半自动

* 打开控制台(cmd)，进入`appman.exe`所在目录:
* 输入 `appman.exe` 并按一次空格
* 拖拽目标程序 (或手动输入路径，注意`"`问题) 到控制台
* 按 `Enter` 完成创建

比如:
```bash
appman.exe "C:\Windows\notepad.exe"
```
没出意外的话，应该创建了一个指向 `appman.exe` 并打开 `notepad.exe` 的快捷方式: `notepad_as_appman.lnk`
否则，查看 `appman.log`
> `nircmd.exe` 的报错未写入日志

### appman.exe

主程序，按需打开外部程序
> 由于主程序是非管理员启动，如果运行的程序需要管理员，会提示权限不够，此时会自动尝试以管理员启动 ( 因为不这么做就打不开程序，至少我不知道解决办法 )

### mksymlink.exe

映射软链接的程序，需要在计划任务中定义

计划任务打开的程序无法传参，故使用临时文件传递参数:

temp_symlink.txt 

> 位于 临时目录\run_app_temp, 如:
> `C:\Users\...\AppData\Local\Temp\run_app_temp`
> 使用完会自动删除

```txt
原始路径1,软链接路径1;原始路径2,软链接路径2
```

### runadmin.exe

以管理员权限打开指定的程序，需要在计划任务中定义

同上，temp_app_runas_admin.txt

```txt
程序1绝对路径 参数
程序2绝对路径 参数
```

### nircmd.exe
[官网](http://nircmd.nirsoft.net)
用来创建快捷方式


## config.json

```json
{
  // notepad.exe会以管理员权限运行，且不会弹出UAC提示
  "nouac": ["C:\\Windows\\notepad.exe"],
  // 需要映射的路径
  "symlink": {
    "r2modman.exe": [
      {
        "dir": "D:\\Game\\r2modman\\App",
        "link": [
          {
            "dirpath": "D:\\Game\\r2modman\\updater",
            "symlink": "%localappdata%\\r2modman-updater"
          },
          {
            "dirpath": "D:\\Game\\r2modman\\App",
            "symlink": "%localappdata%\\Programs\\r2modman"
          }
        ]
      }
    ]
  }
}
```



---



能力有限，代码简陋，欢迎提出更优解决方案
