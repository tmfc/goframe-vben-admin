# 字典管理功能设计文档

## 概述

字典管理功能用于管理系统中的各种固定、可枚举的数据项，如省市区数据、行业分类、学历等级等。这些数据具有以下特点：
- 纯展示数据，不涉及复杂业务逻辑
- 用于选择和展示，不需要特殊处理
- 多个地方复用
- 可能需要多语言支持
- 数据相对稳定，但偶尔会新增

## 适用场景

### ✅ 适合字典管理的场景

| 场景 | 特点 | 说明 |
|------|------|------|
| **省市区数据** | 纯展示、多级联动、偶尔新增 | 不涉及业务逻辑，只需展示和选择 |
| **行业分类** | 纯展示、多语言、可能新增 | 用于企业注册、用户资料等 |
| **学历等级** | 纯展示、多语言、相对稳定 | 用于用户简历、员工档案 |
| **职位类型** | 纯展示、多语言、可能新增 | 用于招聘、员工管理 |
| **工作年限** | 纯展示、相对稳定 | 用于招聘、简历筛选 |
| **性别** | 纯展示、多语言、非常稳定 | 用于用户资料、员工档案 |
| **婚姻状况** | 纯展示、多语言、非常稳定 | 用于用户资料、员工档案 |
| **政治面貌** | 纯展示、相对稳定 | 用于员工档案、人事管理 |
| **血型** | 纯展示、多语言、非常稳定 | 用于员工档案、医疗档案 |
| **证件类型** | 纯展示、多语言、可能新增 | 用于实名认证、用户注册 |

### ❌ 不适合字典管理的场景

- 订单状态（有复杂业务逻辑）
- 用户状态（有复杂业务逻辑）
- 支付方式（需要后端对接、前端页面）
- 审核状态（有复杂业务逻辑）

**核心原则**：字典管理是"数据配置工具"，不是"业务逻辑生成器"。它应该用于已经开发完成的功能的配置，而不是替代代码开发。

## 数据库设计

### 1. 字典类型表 `sys_dict_type`

```sql
CREATE TABLE sys_dict_type (
    id BIGSERIAL PRIMARY KEY,
    type_code VARCHAR(50) NOT NULL UNIQUE COMMENT '字典类型编码，如: payment_method, user_status',
    type_name VARCHAR(100) NOT NULL COMMENT '字典类型名称',
    description VARCHAR(255) COMMENT '描述',
    is_system BOOLEAN DEFAULT FALSE COMMENT '是否系统字典（系统字典不可删除）',
    status SMALLINT DEFAULT 1 COMMENT '状态: 1=启用, 0=禁用',
    sort INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    creator_id BIGINT COMMENT '创建人ID',
    modifier_id BIGINT COMMENT '修改人ID',
    dept_id BIGINT COMMENT '部门ID'
);

CREATE INDEX idx_dict_type_code ON sys_dict_type(type_code);
CREATE INDEX idx_dict_type_status ON sys_dict_type(status);
```

### 2. 字典数据表 `sys_dict_data`

```sql
CREATE TABLE sys_dict_data (
    id BIGSERIAL PRIMARY KEY,
    dict_type_id BIGINT NOT NULL COMMENT '字典类型ID',
    label VARCHAR(100) NOT NULL COMMENT '字典标签（默认语言）',
    value VARCHAR(100) NOT NULL COMMENT '字典值',
    description VARCHAR(255) COMMENT '描述',
    
    -- 多语言支持
    label_en VARCHAR(100) COMMENT '英文标签',
    label_zh_tw VARCHAR(100) COMMENT '繁体中文标签',
    label_ja VARCHAR(100) COMMENT '日文标签',
    
    -- 扩展字段
    color VARCHAR(20) COMMENT '颜色，如: success, error, warning, info',
    icon VARCHAR(100) COMMENT '图标',
    css_class VARCHAR(100) COMMENT 'CSS 类名',
    
    -- 控制字段
    status SMALLINT DEFAULT 1 COMMENT '状态: 1=启用, 0=禁用',
    sort INT DEFAULT 0 COMMENT '排序',
    is_default BOOLEAN DEFAULT FALSE COMMENT '是否默认值',
    
    -- 审计字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    creator_id BIGINT COMMENT '创建人ID',
    modifier_id BIGINT COMMENT '修改人ID',
    dept_id BIGINT COMMENT '部门ID',
    
    CONSTRAINT fk_dict_data_type FOREIGN KEY (dict_type_id) REFERENCES sys_dict_type(id) ON DELETE CASCADE
);

CREATE INDEX idx_dict_data_type_id ON sys_dict_data(dict_type_id);
CREATE INDEX idx_dict_data_value ON sys_dict_data(dict_type_id, value);
CREATE INDEX idx_dict_data_status ON sys_dict_data(status);
```

## 初始化数据

### 1. 性别

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('gender', '性别', '性别选项', TRUE, 1, 1);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, label_en, label_zh_tw, sort) VALUES
(1, '男', '1', '男性', 'Male', '男', 1),
(1, '女', '2', '女性', 'Female', '女', 2),
(1, '未知', '0', '未知性别', 'Unknown', '未知', 3);
```

### 2. 学历等级

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('education', '学历等级', '学历等级列表', TRUE, 1, 2);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, label_en, sort) VALUES
(2, '高中', '1', '高中及以下', 'High School', 1),
(2, '大专', '2', '大专学历', 'College', 2),
(2, '本科', '3', '本科学历', 'Bachelor', 3),
(2, '硕士', '4', '硕士学历', 'Master', 4),
(2, '博士', '5', '博士学历', 'PhD', 5),
(2, '博士后', '6', '博士后', 'Postdoc', 6);
```

### 3. 工作年限

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('work_years', '工作年限', '工作年限选项', TRUE, 1, 3);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, label_en, sort) VALUES
(3, '应届生', '0', '应届毕业生', 'Fresh Graduate', 1),
(3, '1年以下', '1', '1年以下工作经验', 'Under 1 year', 2),
(3, '1-3年', '2', '1-3年工作经验', '1-3 years', 3),
(3, '3-5年', '3', '3-5年工作经验', '3-5 years', 4),
(3, '5-10年', '4', '5-10年工作经验', '5-10 years', 5),
(3, '10年以上', '5', '10年以上工作经验', 'Over 10 years', 6);
```

### 4. 行业分类

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('industry', '行业分类', '企业行业分类', TRUE, 1, 4);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, label_en, sort) VALUES
(4, '互联网/IT', 'internet', '互联网和IT行业', 'Internet/IT', 1),
(4, '金融/银行', 'finance', '金融和银行业', 'Finance/Banking', 2),
(4, '教育/培训', 'education', '教育和培训行业', 'Education/Training', 3),
(4, '医疗/健康', 'healthcare', '医疗和健康行业', 'Healthcare', 4),
(4, '制造业', 'manufacturing', '制造业', 'Manufacturing', 5),
(4, '房地产', 'realestate', '房地产行业', 'Real Estate', 6),
(4, '零售/电商', 'retail', '零售和电商', 'Retail/E-commerce', 7);
```

### 5. 职位类型

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('job_type', '职位类型', '职位类型分类', TRUE, 1, 5);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, label_en, sort) VALUES
(5, '技术类', 'tech', '技术类职位', 'Technology', 1),
(5, '产品类', 'product', '产品类职位', 'Product', 2),
(5, '设计类', 'design', '设计类职位', 'Design', 3),
(5, '运营类', 'operation', '运营类职位', 'Operations', 4),
(5, '市场类', 'marketing', '市场类职位', 'Marketing', 5),
(5, '销售类', 'sales', '销售类职位', 'Sales', 6),
(5, '行政类', 'admin', '行政类职位', 'Administration', 7),
(5, '财务类', 'finance', '财务类职位', 'Finance', 8);
```

### 6. 婚姻状况

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('marital_status', '婚姻状况', '婚姻状况选项', TRUE, 1, 6);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, label_en, sort) VALUES
(6, '未婚', '1', '未婚', 'Single', 1),
(6, '已婚', '2', '已婚', 'Married', 2),
(6, '离异', '3', '离异', 'Divorced', 3),
(6, '丧偶', '4', '丧偶', 'Widowed', 4);
```

### 7. 政治面貌

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('political_status', '政治面貌', '政治面貌选项', TRUE, 1, 7);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, sort) VALUES
(7, '中共党员', '1', '中国共产党党员', 1),
(7, '中共预备党员', '2', '中国共产党预备党员', 2),
(7, '共青团员', '3', '中国共产主义青年团团员', 3),
(7, '群众', '4', '群众', 4),
(7, '民主党派', '5', '民主党派成员', 5),
(7, '无党派人士', '6', '无党派人士', 6);
```

### 8. 血型

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('blood_type', '血型', '血型选项', TRUE, 1, 8);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, label_en, sort) VALUES
(8, 'A型', 'A', 'A型血', 'Type A', 1),
(8, 'B型', 'B', 'B型血', 'Type B', 2),
(8, 'AB型', 'AB', 'AB型血', 'Type AB', 3),
(8, 'O型', 'O', 'O型血', 'Type O', 4),
(8, '未知', 'unknown', '未知血型', 'Unknown', 5);
```

### 9. 证件类型

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES ('id_card_type', '证件类型', '证件类型选项', TRUE, 1, 9);

-- 字典数据
INSERT INTO sys_dict_data (dict_type_id, label, value, description, label_en, sort) VALUES
(9, '身份证', '1', '居民身份证', 'ID Card', 1),
(9, '护照', '2', '护照', 'Passport', 2),
(9, '港澳通行证', '3', '港澳居民来往内地通行证', 'HK/Macau Pass', 3),
(9, '台湾通行证', '4', '台湾居民来往大陆通行证', 'Taiwan Pass', 4),
(9, '军官证', '5', '军官证', 'Military ID', 5),
(9, '其他', '99', '其他证件', 'Other', 6);
```

### 10. 省市区数据

```sql
-- 字典类型
INSERT INTO sys_dict_type (type_code, type_name, description, is_system, status, sort)
VALUES 
('province', '省份', '中国省份列表', TRUE, 1, 10),
('city', '城市', '中国城市列表', TRUE, 1, 11),
('district', '区县', '中国区县列表', TRUE, 1, 12);

-- 字典数据（省份）
INSERT INTO sys_dict_data (dict_type_id, label, value, description, sort) VALUES
(10, '北京市', '110000', '北京市', 1),
(10, '上海市', '310000', '上海市', 2),
(10, '广东省', '440000', '广东省', 3),
(10, '浙江省', '330000', '浙江省', 4);

-- 字典数据（城市）
INSERT INTO sys_dict_data (dict_type_id, label, value, description, sort) VALUES
(11, '北京市', '110100', '北京市', 1),
(11, '上海市', '310100', '上海市', 2),
(11, '广州市', '440100', '广东省广州市', 3),
(11, '深圳市', '440300', '广东省深圳市', 4),
(11, '杭州市', '330100', '浙江省杭州市', 5);

-- 字典数据（区县）
INSERT INTO sys_dict_data (dict_type_id, label, value, description, sort) VALUES
(12, '东城区', '110101', '北京市东城区', 1),
(12, '西城区', '110102', '北京市西城区', 2),
(12, '朝阳区', '110105', '北京市朝阳区', 3),
(12, '海淀区', '110108', '北京市海淀区', 4);
```

## API 设计

### 后端 API

#### 1. 字典类型管理

- `POST /sys-dict-type` - 创建字典类型
- `GET /sys-dict-type/{id}` - 获取字典类型详情
- `PUT /sys-dict-type/{id}` - 更新字典类型
- `DELETE /sys-dict-type/{id}` - 删除字典类型
- `GET /sys-dict-type/list` - 获取字典类型列表

#### 2. 字典数据管理

- `POST /sys-dict-data` - 创建字典数据
- `GET /sys-dict-data/{id}` - 获取字典数据详情
- `PUT /sys-dict-data/{id}` - 更新字典数据
- `DELETE /sys-dict-data/{id}` - 删除字典数据
- `GET /sys-dict-data/list` - 获取字典数据列表

#### 3. 字典查询

- `GET /sys-dict/options/{typeCode}` - 获取字典选项列表（供前端选择使用）
- `GET /sys-dict/label/{typeCode}/{value}` - 获取字典标签（根据值获取标签）

### 前端 API

```typescript
// 字典选项
export interface DictOption {
  label: string;
  value: string;
  color?: string;
  icon?: string;
  isDefault?: boolean;
}

// 字典数据
export interface DictData {
  id: number;
  dictTypeId: number;
  label: string;
  value: string;
  description?: string;
  labelEn?: string;
  labelZhTw?: string;
  labelJa?: string;
  color?: string;
  icon?: string;
  cssClass?: string;
  status: number;
  sort: number;
  isDefault: boolean;
}

// 字典类型
export interface DictType {
  id: number;
  typeCode: string;
  typeName: string;
  description?: string;
  isSystem: boolean;
  status: number;
  sort: number;
}

// API 方法
export async function getDictOptions(typeCode: string): Promise<DictOption[]>;
export async function getDictLabel(typeCode: string, value: string, lang?: string): Promise<string>;
```

## 前端使用示例

### 用户资料表单

```vue
<template>
  <NForm :model="formData">
    <!-- 性别选择 -->
    <NFormItem label="性别">
      <NSelect
        v-model:value="formData.gender"
        :options="genderOptions"
        placeholder="请选择性别"
      />
    </NFormItem>

    <!-- 学历选择 -->
    <NFormItem label="学历">
      <NSelect
        v-model:value="formData.education"
        :options="educationOptions"
        placeholder="请选择学历"
      />
    </NFormItem>

    <!-- 工作年限选择 -->
    <NFormItem label="工作年限">
      <NSelect
        v-model:value="formData.workYears"
        :options="workYearsOptions"
        placeholder="请选择工作年限"
      />
    </NFormItem>

    <!-- 所在城市选择 -->
    <NFormItem label="所在城市">
      <NCascader
        v-model:value="formData.city"
        :options="cityOptions"
        placeholder="请选择城市"
      />
    </NFormItem>
  </NForm>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getDictOptions } from '#/api/sys/dict';

const formData = ref({
  gender: '',
  education: '',
  workYears: '',
  city: [],
});

const genderOptions = ref<DictOption[]>([]);
const educationOptions = ref<DictOption[]>([]);
const workYearsOptions = ref<DictOption[]>([]);
const cityOptions = ref<DictOption[]>([]);

onMounted(async () => {
  // 并行加载所有字典数据
  const [gender, education, workYears, province] = await Promise.all([
    getDictOptions('gender'),
    getDictOptions('education'),
    getDictOptions('work_years'),
    getDictOptions('province'),
  ]);

  genderOptions.value = gender;
  educationOptions.value = education;
  workYearsOptions.value = workYears;
  
  // 构建城市级联数据
  cityOptions.value = await buildCityCascade(province);
});
</script>
```

## 设计要点

| 设计点 | 说明 |
|--------|------|
| **分离类型和数据** | 便于管理和查询，一个类型可以有多个数据项 |
| **多语言支持** | 通过 `label_en`, `label_zh_tw` 等字段支持国际化 |
| **扩展字段** | `color`, `icon`, `css_class` 用于前端展示优化 |
| **状态控制** | `status` 字段控制启用/禁用 |
| **排序支持** | `sort` 字段控制显示顺序 |
| **默认值** | `is_default` 标记默认选项 |
| **系统字典** | `is_system` 标记系统字典，防止误删 |
| **外键约束** | 级联删除，保证数据一致性 |
| **索引优化** | 为常用查询字段添加索引 |

## 注意事项

1. **不要过度使用字典管理**
   - 纯展示、可配置的数据 → 用字典管理
   - 有业务逻辑的状态 → 用常量/枚举

2. **字典管理用于配置型数据**
   - 功能已开发完成，用字典控制显示/隐藏
   - 纯展示数据，不涉及复杂业务逻辑

3. **混合使用**
   - 前端展示用字典
   - 业务逻辑用常量/枚举

4. **添加配置验证**
   - 启动时验证字典配置
   - 防止配置错误导致 bug

5. **缓存策略**
   - 字典数据相对稳定，适合缓存
   - 可以使用 Redis 或内存缓存
   - 缓存失效策略：定时刷新或手动刷新

6. **权限控制**
   - 系统字典不可删除
   - 字典管理功能需要管理员权限