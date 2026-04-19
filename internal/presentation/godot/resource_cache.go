package godot

import "strings"

type ResourceCache struct {
	textures map[string]string
	music    map[string]string
	sfx      map[string]string
}

func NewResourceCache() *ResourceCache {
	return &ResourceCache{
		textures: make(map[string]string),
		music:    make(map[string]string),
		sfx:      make(map[string]string),
	}
}

func (c *ResourceCache) ResolveTexture(path string) string {
	return c.resolve(&c.textures, path)
}

func (c *ResourceCache) ResolveMusic(trackID string) string {
	return c.resolve(&c.music, trackID)
}

func (c *ResourceCache) ResolveSfx(sfxID string) string {
	return c.resolve(&c.sfx, sfxID)
}

func (c *ResourceCache) Stats() (textures int, music int, sfx int) {
	if c == nil {
		return 0, 0, 0
	}
	return len(c.textures), len(c.music), len(c.sfx)
}

func (c *ResourceCache) resolve(target *map[string]string, value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if c == nil {
		return trimmed
	}
	if cached, ok := (*target)[trimmed]; ok {
		return cached
	}
	(*target)[trimmed] = trimmed
	return trimmed
}
