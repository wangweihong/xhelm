# 设计目的
 此项目是一套k8s模板管理仓库，基于helm模板解析功能, 但不包含helm release功能， 即你只能用来管理模板仓库，解析模板，但不能直接部署到kubernetes环境中（helm install not support).


 # 后端
 目前考虑使用etcd存储数据， 并基于etcd进行高可用。
 
 # 流程
 * 添加远程仓库repoA后，如果需要使用repoA的包chartA, 则将chartA先下载到repoA的本地目录中的后，解压到cache目录后进行使用
 * 删除远程仓库repoA时，同时删除本地目录中的数据
 * 远程仓库的chart包不保存到etcd中。

 * 添加chart包到本地目录时，同时会保存到etcd数据库中，并且保存该包的shasum码
 * 每次使用本地仓库的chart包时，会先检测本地目录是否存在该chart包， 存在则进行shasum检验。校验失败则从etcd中重新加载数据。这样保证高可用环境下， 主节点的切换不会导致本地数据的混乱。
 * 每个版本的chart都有一份默认的配置。 用户还可以为每个版本的chart创建多项额外的配置， 额外的配置文件可以是默认配置的子集，只会更新某一部分数据。

 # 考虑
  不保存各仓库的index.yaml文件到etcd, 该文件的大小会随着模板/版本的增多而递增, 如当前kubernetes官方仓库拉下的index.yaml文件已经将近572K, 持续增长将会超过etcd单键能够承载的能力(>1MB).
 # 兼容
  目前不考虑兼容helm命令.
