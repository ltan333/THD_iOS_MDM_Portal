/**
 * Avatar Generator Service
 * Tạo avatar random từ DiceBear API
 * Website: https://www.dicebear.com/
 */

type AvatarStyle = 
  | 'adventurer'
  | 'adventurer-neutral'
  | 'avataaars'
  | 'avataaars-neutral'
  | 'big-ears'
  | 'big-ears-neutral'
  | 'big-smile'
  | 'bottts'
  | 'bottts-neutral'
  | 'croodles'
  | 'croodles-neutral'
  | 'fun-emoji'
  | 'icons'
  | 'identicon'
  | 'initials'
  | 'lorelei'
  | 'micah'
  | 'miniavs'
  | 'notionists'
  | 'notionists-neutral'
  | 'personas'
  | 'pixel-art'
  | 'pixel-art-neutral'
  | 'rings'
  | 'shapes'
  | 'thumbs'
  | 'lorelei'
  | 'domasks'
  | 'la-boulette'
  | 'noto-emoji'
  | 'daboobs'
  | 'open-peeps'
  | 'peerless'
  | 'peeps';

interface AvatarOptions {
  style?: AvatarStyle;
  seed?: string;
  size?: number;
  backgroundColor?: string;
  scale?: number;
  radius?: number;
  translateX?: number;
  translateY?: number;
  rotate?: number;
  flip?: boolean;
}

const DEFAULT_STYLE: AvatarStyle = 'pixel-art';
const DEFAULT_SIZE = 200;
const DICEBEAR_API = 'https://api.dicebear.com/9.x';

/**
 * Tạo avatar URL từ DiceBear API
 * @param seed - Giá trị để generate avatar (ví dụ: email, name, id)
 * @param options - Các option để customize avatar
 * @returns URL của avatar
 * 
 * @example
 * // Sử dụng email làm seed
 * const avatarUrl = generateAvatarUrl('user@example.com');
 * 
 * // Sử dụng custom options
 * const avatarUrl = generateAvatarUrl('John Doe', {
 *   style: 'pixel-art',
 *   size: 256,
 *   backgroundColor: 'FF5733'
 * });
 */
export const generateAvatarUrl = (seed: string, options?: AvatarOptions): string => {
  const {
    style = DEFAULT_STYLE,
    size = DEFAULT_SIZE,
    backgroundColor,
    scale,
    radius,
    translateX,
    translateY,
    rotate,
    flip,
  } = options || {};

  const params = new URLSearchParams();
  
  if (backgroundColor) params.append('backgroundColor', backgroundColor);
  if (scale !== undefined) params.append('scale', String(scale));
  if (radius !== undefined) params.append('radius', String(radius));
  if (translateX !== undefined) params.append('translateX', String(translateX));
  if (translateY !== undefined) params.append('translateY', String(translateY));
  if (rotate !== undefined) params.append('rotate', String(rotate));
  if (flip !== undefined) params.append('flip', String(flip));

  const queryString = params.toString();
  const url = `${DICEBEAR_API}/${style}/${size}/${encodeURIComponent(seed)}.svg`;
  
  return queryString ? `${url}?${queryString}` : url;
};

/**
 * Tạo avatar URL với seed random
 * @param options - Các option để customize avatar
 * @returns URL của avatar với seed random
 * 
 * @example
 * const randomAvatarUrl = generateRandomAvatarUrl();
 */
export const generateRandomAvatarUrl = (options?: AvatarOptions): string => {
  const randomSeed = `user-${Date.now()}-${Math.random().toString(36).substring(7)}`;
  return generateAvatarUrl(randomSeed, options);
};

/**
 * Tạo avatar URL dựa trên tên hoặc email của nhân viên
 * @param firstName - Tên của nhân viên
 * @param lastName - Họ của nhân viên
 * @param options - Các option để customize avatar
 * @returns URL của avatar
 * 
 * @example
 * const avatarUrl = generateEmployeeAvatarUrl('John', 'Doe');
 */
export const generateEmployeeAvatarUrl = (
  firstName: string = '',
  lastName: string = '',
  options?: AvatarOptions
): string => {
  const fullName = `${firstName} ${lastName}`.trim();
  const seed = fullName || `employee-${Date.now()}`;
  
  return generateAvatarUrl(seed, {
    style: 'adventurer-neutral',
    size: 200,
    ...options,
  });
};

/**
 * Tạo avatar từ email
 * @param email - Email của user
 * @param options - Các option để customize avatar
 * @returns URL của avatar
 * 
 * @example
 * const avatarUrl = generateAvatarFromEmail('user@example.com');
 */
export const generateAvatarFromEmail = (email: string, options?: AvatarOptions): string => {
  return generateAvatarUrl(email, {
    style: 'adventurer-neutral',
    size: 200,
    ...options,
  });
};

/**
 * Danh sách các style có sẵn trên DiceBear
 */
export const AVAILABLE_AVATAR_STYLES: AvatarStyle[] = [
  'adventurer',
  'adventurer-neutral',
  'avataaars',
  'avataaars-neutral',
  'big-ears',
  'big-ears-neutral',
  'big-smile',
  'bottts',
  'bottts-neutral',
  'croodles',
  'croodles-neutral',
  'fun-emoji',
  'icons',
  'identicon',
  'initials',
  'lorelei',
  'micah',
  'miniavs',
  'notionists',
  'notionists-neutral',
  'personas',
  'pixel-art',
  'pixel-art-neutral',
  'rings',
  'shapes',
  'thumbs',
  'domasks',
  'la-boulette',
  'noto-emoji',
  'daboobs',
  'open-peeps',
  'peerless',
  'peeps',
];

export default generateAvatarUrl;
