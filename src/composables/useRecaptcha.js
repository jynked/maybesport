import { useReCaptcha } from 'vue-recaptcha-v3'

export const useRecaptcha = () => {
  const { executeRecaptcha, recaptchaLoaded } = useReCaptcha()

  const getToken = async (action = 'submit') => {
    await recaptchaLoaded()
    const token = await executeRecaptcha(action)
    return token
  }

  return { getToken }
}