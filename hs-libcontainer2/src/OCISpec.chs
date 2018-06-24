{-# LANGUAGE ForeignFunctionInterface #-}

#include "oci-spec.h"

module OCISpec (
  LinuxIntelRdt(..),
  l3CacheSchema
) where

import Foreign.Storable (Storable)
import Foreign.C.String (CString)
import Foreign.Ptr (Ptr, nullPtr)
import Control.Monad (liftM)

data LinuxIntelRdt = LinuxIntelRdt {
  l3CacheSchema :: CString
} deriving (Show, Eq)

instance Storable LinuxIntelRdt where
  sizeOf _ = {#sizeof LinuxIntelRdt #}
  alignment _ = {#alignof LinuxIntelRdt #}
  peek ptr = LinuxIntelRdt <$> ({#get LinuxIntelRdt->l3CacheSchema #} ptr)
  poke ptr (LinuxIntelRdt lcs) = do
    {#set LinuxIntelRdt.l3CacheSchema #} ptr lcs