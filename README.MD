#### filemd5

Android应用维护升级共计

#### 工具说明
获取当前目录下的 .apk 文件的更新时间和MD5值
并将获取到的值更新到下级目录内的 .json 文件中
根据 name 对应：
 - 将 version 复制 .apk 文件的更新时间
 - 将 hash 复制 .apk 文件的MD5值

#### 工具参数
- exclude_dir 过滤目录名称，多个目录用 ，隔开； 默认过滤 com.chindeo.launcher.app
- exclude_file 过滤apk文件名称，多个文件用 ，隔开

过滤 com.chindeo.launcher.app 目录