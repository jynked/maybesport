import { createI18n } from 'vue-i18n';

const messages = {
    en: {
        catalogHeaderNav: 'Catalog',
        altLogo: 'Logo',
        cartHeaderNav: 'Cart',
        lkHeaderNav: 'LK',
        favouritesHeaderNav: 'Favourites',
        contactsHeaderNav: 'Contacts',
        ordersHeaderNav: 'Orders',
        altLK: 'Personal Account',
        altCart: 'Cart',
        newItemImageAlt: 'Image of new item',
        newItemMore: 'More',
        itemImageAlt: 'Item image',
        catalog: 'Catalog',
        lookForCatalog: 'Take a look at our catalog',
        color: "Color",
        from: "From",
        available: "Available",
        on_request: "On request",
        out_of_stock: "Out of stock",
        ourTGC: "Our Telegram Channel",
        TGCBGAlt: "Telegram Channel's image",
        TGAlt: "Telegram image",
        TGC: "TGC",
        link: "link",
        privacyPolicy: 'Privacy policy',
        userAgreement: 'User agreement'
    },
    ru: {
        catalogHeaderNav: 'Каталог',
        altLogo: 'Логотип',
        cartHeaderNav: 'Корзина',
        lkHeaderNav: 'ЛК',
        favouritesHeaderNav: 'Избранное',
        contactsHeaderNav: 'Контакты',
        ordersHeaderNav: 'Заказы',
        altLK: 'Личный кабинет',
        altCart: 'Корзина',
        newItemImageAlt: 'Изображение нового товара',
        newItemMore: 'Перейти',
        itemImageAlt: 'Изображение товара',
        catalog: 'Каталог',
        lookForCatalog: 'Загляните к нам в каталог',
        color: "Цвет",
        from: "От",
        available: "В наличии",
        on_request: "Под заказ",
        out_of_stock: "Нет в наличии",
        ourTGC: "Наш Телеграмм канал",
        TGCBGAlt: "Изображение телеграм канала",
        TGAlt: "Изображение телеграм",
        TGC: "ТГК",
        link: "линк",
        privacyPolicy: 'Политика конфиденциальности',
        userAgreement: 'Пользовательское соглашение'
    },
};

const i18n = createI18n({
    locale: 'ru',
    messages,
});

export default i18n;