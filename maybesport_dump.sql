--
-- PostgreSQL database dump
--

\restrict QDgw2H71m3VLGnAInNiZRhsivAAoiGpdW0ApNbyxeTi91PPLmpfqfhro8kmbQbb

-- Dumped from database version 18.3
-- Dumped by pg_dump version 18.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: main_page_new; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.main_page_new (
    id integer NOT NULL,
    title_ru text,
    title_en text,
    image text,
    product_item_unique_id text
);


ALTER TABLE public.main_page_new OWNER TO postgres;

--
-- Name: main_page_new_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.main_page_new_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.main_page_new_id_seq OWNER TO postgres;

--
-- Name: main_page_new_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.main_page_new_id_seq OWNED BY public.main_page_new.id;


--
-- Name: order_item_status_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_item_status_history (
    id integer NOT NULL,
    order_item_id integer NOT NULL,
    status text NOT NULL,
    description text,
    "timestamp" timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.order_item_status_history OWNER TO postgres;

--
-- Name: order_item_status_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.order_item_status_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_item_status_history_id_seq OWNER TO postgres;

--
-- Name: order_item_status_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_item_status_history_id_seq OWNED BY public.order_item_status_history.id;


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_items (
    id integer NOT NULL,
    order_id integer,
    product_item_id integer,
    size text,
    price bigint NOT NULL,
    quantity bigint NOT NULL,
    status text DEFAULT 'created'::text
);


ALTER TABLE public.order_items OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.order_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_items_id_seq OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- Name: order_status_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_status_history (
    id integer NOT NULL,
    order_id integer,
    status text NOT NULL,
    description text,
    "timestamp" timestamp with time zone DEFAULT now()
);


ALTER TABLE public.order_status_history OWNER TO postgres;

--
-- Name: order_status_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.order_status_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_status_history_id_seq OWNER TO postgres;

--
-- Name: order_status_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_status_history_id_seq OWNED BY public.order_status_history.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    user_id integer,
    status text DEFAULT 'created'::text NOT NULL,
    total_amount bigint NOT NULL,
    delivery_address text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: product_item_sizes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product_item_sizes (
    id integer NOT NULL,
    product_item_id integer,
    size text NOT NULL,
    price bigint NOT NULL,
    is_on_request boolean DEFAULT false,
    quantity bigint DEFAULT 0 NOT NULL,
    price_cny bigint
);


ALTER TABLE public.product_item_sizes OWNER TO postgres;

--
-- Name: product_item_sizes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.product_item_sizes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_item_sizes_id_seq OWNER TO postgres;

--
-- Name: product_item_sizes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.product_item_sizes_id_seq OWNED BY public.product_item_sizes.id;


--
-- Name: product_item_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product_item_tags (
    product_item_id integer NOT NULL,
    tag_id integer NOT NULL
);


ALTER TABLE public.product_item_tags OWNER TO postgres;

--
-- Name: product_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product_items (
    id integer NOT NULL,
    product_id integer,
    unique_id text NOT NULL,
    images text[],
    color_ru text,
    color_en text
);


ALTER TABLE public.product_items OWNER TO postgres;

--
-- Name: product_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.product_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_items_id_seq OWNER TO postgres;

--
-- Name: product_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.product_items_id_seq OWNED BY public.product_items.id;


--
-- Name: product_structures; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product_structures (
    id integer NOT NULL,
    product_id integer,
    name_ru text,
    name_en text,
    percent integer
);


ALTER TABLE public.product_structures OWNER TO postgres;

--
-- Name: product_structures_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.product_structures_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_structures_id_seq OWNER TO postgres;

--
-- Name: product_structures_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.product_structures_id_seq OWNED BY public.product_structures.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id integer NOT NULL,
    type_ru text,
    type_en text,
    title_ru text,
    title_en text,
    desc_ru text,
    desc_en text,
    brand text,
    country_ru text,
    country_en text,
    category_ru text,
    category_en text,
    created_at timestamp with time zone
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tags (
    id integer NOT NULL,
    name_ru text,
    name_en text
);


ALTER TABLE public.tags OWNER TO postgres;

--
-- Name: tags_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tags_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tags_id_seq OWNER TO postgres;

--
-- Name: tags_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tags_id_seq OWNED BY public.tags.id;


--
-- Name: user_cart; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_cart (
    id integer NOT NULL,
    user_id integer,
    product_item_id integer,
    size text,
    quantity bigint DEFAULT 1 NOT NULL,
    added_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.user_cart OWNER TO postgres;

--
-- Name: user_cart_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_cart_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_cart_id_seq OWNER TO postgres;

--
-- Name: user_cart_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_cart_id_seq OWNED BY public.user_cart.id;


--
-- Name: user_favourites; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_favourites (
    id integer NOT NULL,
    user_id integer,
    product_item_id integer,
    size text,
    added_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.user_favourites OWNER TO postgres;

--
-- Name: user_favourites_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_favourites_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_favourites_id_seq OWNER TO postgres;

--
-- Name: user_favourites_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_favourites_id_seq OWNED BY public.user_favourites.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    name text,
    is_admin boolean DEFAULT false
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: main_page_new id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.main_page_new ALTER COLUMN id SET DEFAULT nextval('public.main_page_new_id_seq'::regclass);


--
-- Name: order_item_status_history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_item_status_history ALTER COLUMN id SET DEFAULT nextval('public.order_item_status_history_id_seq'::regclass);


--
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- Name: order_status_history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_status_history ALTER COLUMN id SET DEFAULT nextval('public.order_status_history_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: product_item_sizes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_item_sizes ALTER COLUMN id SET DEFAULT nextval('public.product_item_sizes_id_seq'::regclass);


--
-- Name: product_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_items ALTER COLUMN id SET DEFAULT nextval('public.product_items_id_seq'::regclass);


--
-- Name: product_structures id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_structures ALTER COLUMN id SET DEFAULT nextval('public.product_structures_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: tags id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags ALTER COLUMN id SET DEFAULT nextval('public.tags_id_seq'::regclass);


--
-- Name: user_cart id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cart ALTER COLUMN id SET DEFAULT nextval('public.user_cart_id_seq'::regclass);


--
-- Name: user_favourites id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_favourites ALTER COLUMN id SET DEFAULT nextval('public.user_favourites_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: main_page_new; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.main_page_new (id, title_ru, title_en, image, product_item_unique_id) FROM stdin;
1	Новинка недели - Adidas Predator в синем цвете!	The new product of the week is the Adidas Predator in blue!	https://avatars.dzeninfra.ru/get-zen_doc/1654945/pub_5e4f865bf2b93d016c119250_5e4f988d11638a2a18c12dbe/orig	2-1
\.


--
-- Data for Name: order_item_status_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_item_status_history (id, order_item_id, status, description, "timestamp") FROM stdin;
60	51	processing		2026-04-21 15:49:23.132291+03
61	52	created	Заказ оформлен	2026-04-21 15:50:18.080884+03
62	53	created	Заказ оформлен	2026-04-21 15:50:18.080884+03
63	52	received		2026-04-21 15:50:43.99059+03
64	52	processing		2026-04-21 15:50:56.191192+03
65	52	cancelled		2026-04-21 15:51:03.902008+03
53	45	created	Заказ оформлен	2026-04-21 14:08:58.451464+03
54	46	created	Заказ оформлен	2026-04-21 14:08:58.451464+03
55	47	created	Заказ оформлен	2026-04-21 14:31:56.697464+03
56	48	created	Заказ оформлен	2026-04-21 14:31:56.697464+03
57	49	created	Заказ оформлен	2026-04-21 15:26:43.808547+03
58	50	created	Заказ оформлен	2026-04-21 15:47:25.512512+03
59	51	created	Заказ оформлен	2026-04-21 15:47:25.512512+03
\.


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_items (id, order_id, product_item_id, size, price, quantity, status) FROM stdin;
48	21	5	M	2377	1	created
49	22	5	L	1783	4	created
52	24	6	XL	264142	1	cancelled
45	20	\N	39	14317431	12	created
46	20	\N	41	257206666	122	created
47	21	\N	41	257206666	161	created
50	23	\N	45	1300	5	created
51	23	\N	44	2576	8	processing
53	24	\N	40	1462536	7	created
\.


--
-- Data for Name: order_status_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_status_history (id, order_id, status, description, "timestamp") FROM stdin;
67	20	created	Заказ оформлен	2026-04-21 14:08:58.451464+03
68	21	created	Заказ оформлен	2026-04-21 14:31:56.697464+03
69	22	created	Заказ оформлен	2026-04-21 15:26:43.808547+03
70	23	created	Заказ оформлен	2026-04-21 15:47:25.512512+03
71	24	created	Заказ оформлен	2026-04-21 15:50:18.080884+03
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, user_id, status, total_amount, delivery_address, created_at, updated_at) FROM stdin;
20	10	created	31551022424	мой дом епта	2026-04-21 14:08:58.451464+03	2026-04-21 14:08:58.451464+03
21	10	created	41410275603	sdfsd fsdf sd	2026-04-21 14:31:56.697464+03	2026-04-21 14:31:56.697464+03
22	10	created	7132	sdfdsf dsfsdf s	2026-04-21 15:26:43.808547+03	2026-04-21 15:26:43.808547+03
23	10	processing	27108	jhhkjhkjhjkh	2026-04-21 15:47:25.512512+03	2026-04-21 15:49:23.132291+03
24	7	created	10501894	lkjkljlk	2026-04-21 15:50:18.080884+03	2026-04-21 15:51:03.902008+03
\.


--
-- Data for Name: product_item_sizes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.product_item_sizes (id, product_item_id, size, price, is_on_request, quantity, price_cny) FROM stdin;
11	5	M	2377	f	0	200
13	6	S	2639	t	0	222
15	7	UNI	1189	f	20	100
43	25	41	104	t	0	9
12	5	L	1783	f	0	150
14	6	XL	264142	f	1	22222
44	26	41	257206666	f	583	22365797
45	26	45	1300	f	4	113
46	26	44	2576	t	102	224
47	26	123	375843	t	122	32682
48	26	31	2701351599	t	220	234900139
49	26	39	14317431	f	85	1244994
50	27	40	1462536	t	33326	127177
9	3	43	1462	f	0	123
8	3	42	2784222	f	0	234234
10	4	41	132071	f	0	11111
\.


--
-- Data for Name: product_item_tags; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.product_item_tags (product_item_id, tag_id) FROM stdin;
3	1
3	2
3	6
4	7
5	8
5	2
25	1
25	6
25	2
25	8
25	7
26	1
26	2
27	2
\.


--
-- Data for Name: product_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.product_items (id, product_id, unique_id, images, color_ru, color_en) FROM stdin;
3	2	2-1	{https://img.freepik.com/premium-vector/sport-shoes-flat-icon-fitness-basketball-gym-sign-vector-illustration_939711-3427.jpg?semt=ais_hybrid&w=740}	Черный	Black
4	2	2-2	{https://i.pinimg.com/originals/26/5e/29/265e29a56676d3b4e13195ca0b5940ca.jpg}	Красный	Red
5	3	3-1	{https://img.freepik.com/premium-vector/vector-colorful-shoe-with-white-laces-vibrant-background_410516-100327.jpg?semt=ais_hybrid&w=740}	Темно-синий	Dark Blue
6	3	3-2	{https://avatars.mds.yandex.net/i?id=e6ac10f13a3383a696f3a4e4c75000b28bb784f1-16320731-images-thumbs&n=13}	Черный	Black
7	4	4-1	{https://i.pinimg.com/originals/9e/c2/69/9ec269a2c8ab04b53ca21acb1df3f485.png}	Серый	Gray
25	5	5-1	{/uploads/a8356b34-09c6-4d94-ada3-3a569dc2bd0b.png,/uploads/aa089897-b0e9-4158-bcba-bcc11468ce6c.png,/uploads/ffc56977-236e-42b6-88ee-debd500bf0e2.png}	Красный	Red
26	1	1-1	{https://i.pinimg.com/originals/9e/c2/69/9ec269a2c8ab04b53ca21acb1df3f485.png,https://i.pinimg.com/originals/9e/c2/69/9ec269a2c8ab04b53ca21acb1df3f485.png,https://avatars.mds.yandex.net/i?id=8fa051f1be50f98a842f1a312639505b_l-5255597-images-thumbs&n=13,https://i.pinimg.com/originals/9e/c2/69/9ec269a2c8ab04b53ca21acb1df3f485.png,https://i.pinimg.com/originals/9e/c2/69/9ec269a2c8ab04b53ca21acb1df3f485.png,https://avatars.mds.yandex.net/i?id=8fa051f1be50f98a842f1a312639505b_l-5255597-images-thumbs&n=13,https://i.pinimg.com/originals/9e/c2/69/9ec269a2c8ab04b53ca21acb1df3f485.png,https://i.pinimg.com/originals/9e/c2/69/9ec269a2c8ab04b53ca21acb1df3f485.png,https://avatars.mds.yandex.net/i?id=8fa051f1be50f98a842f1a312639505b_l-5255597-images-thumbs&n=13}	Синий	Blue
27	1	1-2	{https://i.pinimg.com/originals/3e/fd/13/3efd133b76a29919b088ed16776b5c0a.png,/uploads/05ea1bfb-12d2-41cc-9613-70071323c1e7.png}	Желтый	Yellow
\.


--
-- Data for Name: product_structures; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.product_structures (id, product_id, name_ru, name_en, percent) FROM stdin;
3	2	Синтетика	Synthetic	70
4	2	Резина	Rubber	15
5	3	Синтетика	Synthetic	70
6	3	Резина	Rubber	15
7	4	Синтетика	Synthetic	100
12	1	Текстиль	Textile	70
13	1	Резина	Rubber	30
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, type_ru, type_en, title_ru, title_en, desc_ru, desc_en, brand, country_ru, country_en, category_ru, category_en, created_at) FROM stdin;
2	Бутсы	Cleats	Бутсы Nike Mercurial	Nike Mercurial Boots	Легкие и скоростные бутсы для наилучшего сцепления и маневренности на поле.	Lightweight and speed boots for the best traction and agility on the field.	Nike	США	USA	Обувь	Footwear	2023-09-15 03:00:00+03
3	Пуховики	Down-jackets	Пуховик The North Face	The North Face Down Jacket	Теплый и легкий пуховик для экстремально низких температур. Идеален для зимних видов спорта.	Warm and lightweight down jacket for extremely low temperatures. Perfect for winter sports.	The North Face	США	USA	Верхняя одежда	Outerwear	2023-11-05 03:00:00+03
4	Балаклавы	Balaclavas	Балаклава Craft	Craft Balaclava	Теплая и дышащая балаклава для бега и лыжного спорта зимой.	Warm and breathable balaclava for running and skiing in winter.	Craft	Швеция	Sweden	Аксессуары	Accessories	2023-12-10 03:00:00+03
5	Бутсы	Cleats	Бутсы NIKE Zoom Mercurial Vapor 16 Elite FG	Бутсы Nike Mercurial Vapor 16 Academ	Бутсы выполнены из синтетической кожи NikeSkin с текстильной подкладкой и подошвой из ТПУ, модель Mercurial Vapor 16 Academy. Детали: верх из технологии NikeSkin для улучшенного контроля мяча, система шнуровки, вставка Zoom Air для амортизации и отзывчивости, внутренний каркас для поддержки и стабильности стопы, низкопрофильная конструкция для лучшей стабильности, гибридные шипы для максимального сцепления с натуральным и искусственным газоном, петли на язычке и заднике для удобства надевания. Рекомендованы для игры на натуральных полях и современном искусственном газоне, температурный режим эксплуатации выше +5°C.		Nike	Вьетнам	Vietnam	Обувь	Footwear	2026-03-04 22:05:30.312+03
1	Кроссовки	Sneakers	Кроссовки Adidas Predator	Sneakers Adidas Predator	Футбольные кроссовки для искусственного покрытия с технологией контроля мяча.	Football sneakers for artificial turf with ball control technology.	Adidas	Германия	Germany	Обувь	Footwear	2023-10-01 03:00:00+03
\.


--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tags (id, name_ru, name_en) FROM stdin;
6	Новое	New
8	Идея для подарка	Gift idea
7	Ограничено	Limited
1	Хит	Hit
2	Популярное	Popular
\.


--
-- Data for Name: user_cart; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_cart (id, user_id, product_item_id, size, quantity, added_at, updated_at) FROM stdin;
116	10	3	43	5	2026-04-21 15:46:40.755558+03	2026-04-21 15:46:40.755558+03
112	7	4	41	2	2026-04-21 15:35:17.227439+03	2026-04-21 15:35:17.227439+03
\.


--
-- Data for Name: user_favourites; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_favourites (id, user_id, product_item_id, size, added_at) FROM stdin;
39	7	4	41	2026-04-21 15:34:56.82793+03
40	7	6	XL	2026-04-21 15:35:13.600537+03
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, password, created_at, name, is_admin) FROM stdin;
10	qwe@qwe.qwe	$2a$10$cULwhwqM4XEeDT9GT9Y6Du6tRKXBdZM1wu/qRf8JzTJL/ppqKUnL.	2026-04-21 09:22:39.77867+03		f
7	alexeyandboryaandmaxformaybesport222_sdf_idkbtw@gmail.com	$2a$10$B0IPvTFgcnH/nFpPWsRPuOW6k0uFFGpxNgzIM2ojbQWPMW5cKsyx2	2026-04-21 05:30:34.967098+03	Айти отдел	t
\.


--
-- Name: main_page_new_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.main_page_new_id_seq', 1, false);


--
-- Name: order_item_status_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_item_status_history_id_seq', 65, true);


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_items_id_seq', 53, true);


--
-- Name: order_status_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_status_history_id_seq', 71, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.orders_id_seq', 24, true);


--
-- Name: product_item_sizes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.product_item_sizes_id_seq', 50, true);


--
-- Name: product_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.product_items_id_seq', 27, true);


--
-- Name: product_structures_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.product_structures_id_seq', 13, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 1, false);


--
-- Name: tags_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tags_id_seq', 88, true);


--
-- Name: user_cart_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_cart_id_seq', 118, true);


--
-- Name: user_favourites_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_favourites_id_seq', 43, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 10, true);


--
-- Name: main_page_new main_page_new_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.main_page_new
    ADD CONSTRAINT main_page_new_pkey PRIMARY KEY (id);


--
-- Name: order_item_status_history order_item_status_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_item_status_history
    ADD CONSTRAINT order_item_status_history_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: order_status_history order_status_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_status_history
    ADD CONSTRAINT order_status_history_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: product_item_sizes product_item_sizes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_item_sizes
    ADD CONSTRAINT product_item_sizes_pkey PRIMARY KEY (id);


--
-- Name: product_item_tags product_item_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_item_tags
    ADD CONSTRAINT product_item_tags_pkey PRIMARY KEY (product_item_id, tag_id);


--
-- Name: product_items product_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_items
    ADD CONSTRAINT product_items_pkey PRIMARY KEY (id);


--
-- Name: product_items product_items_unique_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_items
    ADD CONSTRAINT product_items_unique_id_key UNIQUE (unique_id);


--
-- Name: product_structures product_structures_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_structures
    ADD CONSTRAINT product_structures_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: tags tags_name_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_name_unique UNIQUE (name_ru, name_en);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: user_cart user_cart_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cart
    ADD CONSTRAINT user_cart_pkey PRIMARY KEY (id);


--
-- Name: user_cart user_cart_user_id_product_item_id_size_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cart
    ADD CONSTRAINT user_cart_user_id_product_item_id_size_key UNIQUE (user_id, product_item_id, size);


--
-- Name: user_favourites user_favourites_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_favourites
    ADD CONSTRAINT user_favourites_pkey PRIMARY KEY (id);


--
-- Name: user_favourites user_favourites_user_id_product_item_id_size_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_favourites
    ADD CONSTRAINT user_favourites_user_id_product_item_id_size_key UNIQUE (user_id, product_item_id, size);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_order_item_status_history_order_item_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_order_item_status_history_order_item_id ON public.order_item_status_history USING btree (order_item_id);


--
-- Name: idx_order_items_order_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_order_items_order_id ON public.order_items USING btree (order_id);


--
-- Name: idx_orders_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_orders_user_id ON public.orders USING btree (user_id);


--
-- Name: idx_product_item_sizes_product_item_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_product_item_sizes_product_item_id ON public.product_item_sizes USING btree (product_item_id);


--
-- Name: idx_product_items_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_product_items_product_id ON public.product_items USING btree (product_id);


--
-- Name: idx_product_items_unique_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_product_items_unique_id ON public.product_items USING btree (unique_id);


--
-- Name: idx_product_structures_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_product_structures_product_id ON public.product_structures USING btree (product_id);


--
-- Name: idx_user_cart_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_cart_user_id ON public.user_cart USING btree (user_id);


--
-- Name: idx_user_favourites_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_favourites_user_id ON public.user_favourites USING btree (user_id);


--
-- Name: main_page_new main_page_new_product_item_unique_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.main_page_new
    ADD CONSTRAINT main_page_new_product_item_unique_id_fkey FOREIGN KEY (product_item_unique_id) REFERENCES public.product_items(unique_id) ON DELETE SET NULL;


--
-- Name: order_item_status_history order_item_status_history_order_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_item_status_history
    ADD CONSTRAINT order_item_status_history_order_item_id_fkey FOREIGN KEY (order_item_id) REFERENCES public.order_items(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_product_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_item_id_fkey FOREIGN KEY (product_item_id) REFERENCES public.product_items(id) ON DELETE SET NULL;


--
-- Name: order_status_history order_status_history_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_status_history
    ADD CONSTRAINT order_status_history_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: product_item_sizes product_item_sizes_product_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_item_sizes
    ADD CONSTRAINT product_item_sizes_product_item_id_fkey FOREIGN KEY (product_item_id) REFERENCES public.product_items(id) ON DELETE CASCADE;


--
-- Name: product_item_tags product_item_tags_product_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_item_tags
    ADD CONSTRAINT product_item_tags_product_item_id_fkey FOREIGN KEY (product_item_id) REFERENCES public.product_items(id) ON DELETE CASCADE;


--
-- Name: product_item_tags product_item_tags_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_item_tags
    ADD CONSTRAINT product_item_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id) ON DELETE CASCADE;


--
-- Name: product_items product_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_items
    ADD CONSTRAINT product_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: product_structures product_structures_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_structures
    ADD CONSTRAINT product_structures_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: user_cart user_cart_product_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cart
    ADD CONSTRAINT user_cart_product_item_id_fkey FOREIGN KEY (product_item_id) REFERENCES public.product_items(id) ON DELETE CASCADE;


--
-- Name: user_cart user_cart_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_cart
    ADD CONSTRAINT user_cart_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: user_favourites user_favourites_product_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_favourites
    ADD CONSTRAINT user_favourites_product_item_id_fkey FOREIGN KEY (product_item_id) REFERENCES public.product_items(id) ON DELETE CASCADE;


--
-- Name: user_favourites user_favourites_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_favourites
    ADD CONSTRAINT user_favourites_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict QDgw2H71m3VLGnAInNiZRhsivAAoiGpdW0ApNbyxeTi91PPLmpfqfhro8kmbQbb

